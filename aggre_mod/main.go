package main

import (
	"os"
    "net"
	"net/http"
    "flag"
    "bufio"
	"log"
	"io"
	"io/ioutil"
	"time"
	"strconv"

	"github.com/gorilla/handlers"
    "github.com/najeira/ltsv"
	"github.com/fluent/fluent-logger-golang/fluent"
)

var (
	addr = flag.String("addr", "0.0.0.0:8000", "The address to bind the server to.")
	fluentdHost = flag.String("fluentd-host", "localhost", "The fluentd host.")
	fluentdPort = flag.Int("fluentd-port", 24224, "The fluentd port.")
	timeout = flag.Int("timeout", 30, "The device read timeout in seconds.")
	debug = flag.Bool("debug", false, "Enable debug logging.")
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type Device struct {
	Name string
	Address string
}

type DataPoint struct {
	Device Device
	Data map[string]string
}

// Channel used to send data to the processData() goroutine
var dataChan = make(chan DataPoint, 1000)

// Channel used to 
var deviceChan = make(chan Device, 100)

func initLogging() {
	debugOut := ioutil.Discard
	if *debug {
		debugOut = os.Stdout
	}
	Debug = log.New(debugOut, "[DEBUG] ", log.Ldate|log.Ltime)
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	Warning = log.New(os.Stderr, "[WARNING] ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
}

type DeviceHandler struct {
	Device Device
	Done bool
}

var deviceHandlers = make(map[string]*DeviceHandler)

// Connects to a device and receives data.
func (dh *DeviceHandler) Handle() {
	defer dh.Finish()

	Info.Printf("Connecting to %s...", dh.Device.Address)

	c, err := net.DialTimeout("tcp", dh.Device.Address, time.Duration(*timeout) * time.Second)
	if err != nil {
		Error.Printf("Could not connect to %s: %v", dh.Device.Address, err)
		return
    }

	Info.Printf("Connected to %s", dh.Device.Address)

	defer c.Close()

    reader := ltsv.NewReader(bufio.NewReader(c))

	// Start reading data.
	for {
		// Set the read timeout.
		c.SetDeadline(time.Now().Add(time.Duration(*timeout) * time.Second))
        if data, err := reader.Read(); err == nil {
			// Send the received data to the data channel.
			dataChan <- DataPoint{
				Device: dh.Device,
				Data: data,
			}
        } else {
			if err == io.EOF {
				Warning.Printf("Connection to %s terminated", dh.Device.Address)
				return
			} else if e,ok := err.(net.Error); ok && e.Timeout() {
				Warning.Printf("Connection to %s has timed out", dh.Device.Address)
				return
			}

			Error.Printf("Error reading data from %s (%s): %v", dh.Device.Name, dh.Device.Address, err)
		}
    }
}

func (dh *DeviceHandler) Finish() {
	dh.Done = true
}

func cleanupDevices() {
	// Create a goroutine to clean up device handlers.
	for {
		// Delete finished device handlers.
		var toDelete []string
		for n,h := range deviceHandlers {
			if h.Done {
				Info.Printf("%s disconnected.", n)
				toDelete = append(toDelete, n)
			}
		}
		for _,n := range toDelete {
			Info.Printf("Removing %s.", n)
			delete(deviceHandlers, n)
		}

		// Only check every 5 seconds so we don't use too much CPU.
		time.Sleep(5 * time.Second)
	}
}

func addFloatValue(name string, jsonValue map[string]interface{}, data DataPoint) {
	if data.Data[name] != "" {
		if val, err := strconv.ParseFloat(data.Data[name], 64); err == nil {
			jsonValue[name] = val
		} else {
			Error.Printf("Error parsing %s data from %s: %v", name, data.Device.Name, err)
		}
	}
}

// Proceses data in parallel
func processData(logger *fluent.Fluent) {

	for {
		// Block on the data channel.
		data := <-dataChan

		jsonValue := make(map[string]interface{})

		jsonValue["name"] = data.Device.Name

		timestamp, err := strconv.ParseInt(data.Data["timestamp"], 10, 64)
		if err != nil {
			Error.Printf("Error reading timestamp from %s (%s): %v", data.Device.Name, data.Device.Address, err)
			continue
		}

		jsonValue["timestamp"] = timestamp
		addFloatValue("temp", jsonValue, data)
		addFloatValue("humidity", jsonValue, data)
		addFloatValue("windspeed", jsonValue, data)
		addFloatValue("winddirection", jsonValue, data)
		addFloatValue("rainfall", jsonValue, data)

		// Send data directly to Fluentd 
		logger.Post("aggre_mod.sensordata", jsonValue)

		if err != nil {
			Error.Printf("Could not send data from %s to Fluentd: %v", data.Device.Name, err)
		} else {
			Debug.Printf("Data processed (%s): %s", data.Device.Name, data.Data)
		}
	}
}

// A simple io.Writer wrapper around a logger so that we can use
// the logger as an io.Writer
type LogWriter struct { *log.Logger }
func (w LogWriter) Write(b []byte) (int, error) {
      w.Printf("%s", b)
      return len(b), nil
}

func apiServer() {
    http.Handle("/api/devices", handlers.CombinedLoggingHandler(LogWriter{Info}, handlers.MethodHandler{
        "POST": http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			name := r.PostFormValue("name")
			address := r.PostFormValue("address")

			if (name == "" || address == "") {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			// TODO: Check if there is a duplicate name.
			d := Device{
				Name: name,
				Address: address,
			}

			// Add new device handlers.
			Info.Printf("Got new device %s at %s", d.Name, d.Address)

			deviceHandlers[d.Name] = &DeviceHandler{
				Device: d,
				Done: false,
			}

			// Create a goroutine to get data.
			go deviceHandlers[d.Name].Handle()
        }),
    }))

    Info.Printf("Listening on %s...", *addr)
    Error.Fatal(http.ListenAndServe(*addr, nil))
}

func main() {
    flag.Parse()

	initLogging()

	logger, err := fluent.New(fluent.Config{
		FluentHost: *fluentdHost,
		FluentPort: *fluentdPort,
	})
	if err != nil {
		Error.Fatal("Could not connect to Fluentd.", err)
	}

	// Start the thead to process data.
	go processData(logger)

	// Start the api server
	go apiServer()

	// Run a routine to clean up disconnected devices in the main thread.
	cleanupDevices()
}
