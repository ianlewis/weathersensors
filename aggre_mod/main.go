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

	"github.com/gorilla/handlers"
    "github.com/najeira/ltsv"
)

var (
	addr = flag.String("addr", "0.0.0.0:8000", "The address to bind the server to.")
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

var dataChan = make(chan DataPoint, 1000)
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

			Error.Printf("Error reading data from %s: %v", dh.Device.Address, err)
		}
    }
}

func (dh *DeviceHandler) Finish() {
	dh.Done = true
}

func monitorDevices() {
	for {
		select {
		case d := <-deviceChan:
			// Add new device handlers.
			Info.Printf("Got new device %s at %s", d.Name, d.Address)

			deviceHandlers[d.Name] = &DeviceHandler{
				Device: d,
				Done: false,
			}
			go deviceHandlers[d.Name].Handle()
		default:
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
		}
	}
}

// Proceses data in parallel
func processData() {
	for {
		// Block on the data channel.
		data := <-dataChan
		// Got a data point
		// TODO: Send data to Fluentd (or directly to BigQuery?)
		Debug.Printf("%s: %s", data.Device.Name, data.Data)
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
			deviceChan <- Device{
				Name: name,
				Address: address,
			}
        }),
    }))

    Info.Printf("Listening on %s...", *addr)
    Error.Fatal(http.ListenAndServe(*addr, nil))
}


func main() {
    flag.Parse()

	initLogging()

	// Start the thead to process data.
	go processData()

	// Start the api server
	go apiServer()

	// Start the monitor
	monitorDevices()
}
