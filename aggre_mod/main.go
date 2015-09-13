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

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"github.com/google/google-api-go-client/bigquery/v2"  // bigquery
)

var (
	addr = flag.String("addr", "0.0.0.0:8000", "The address to bind the server to.")
	bqproject = flag.String("project", "", "The BigQuery project.")
	bqdataset = flag.String("dataset", "", "The BigQuery dataset.")
	bqtable = flag.String("table", "", "The BigQuery table to stream to.")
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

			Error.Printf("Error reading data from %s (%s): %v", dh.Device.Name, dh.Device.Address, err)
		}
    }
}

func (dh *DeviceHandler) Finish() {
	dh.Done = true
}

func monitorDevices() {
	// Create a goroutine to clean up device handlers.
	go func() {
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

			time.Sleep(5 * time.Second)
		}
	}()

	for {
		// Block on the device channel.
		d := <-deviceChan

		// Add new device handlers.
		Info.Printf("Got new device %s at %s", d.Name, d.Address)

		deviceHandlers[d.Name] = &DeviceHandler{
			Device: d,
			Done: false,
		}
		go deviceHandlers[d.Name].Handle()
	}
}

func addFloatValue(name string, jsonValue map[string]bigquery.JsonValue, data DataPoint) {
	if data.Data[name] != "" {
		if val, err := strconv.ParseFloat(data.Data[name], 64); err == nil {
			jsonValue[name] = val
		} else {
			Error.Printf("Error parsing %s data from %s: %v", name, data.Device.Name, err)
		}
	}
}

// Proceses data in parallel
func processData(s *bigquery.Service) {

	for {
		// Block on the data channel.
		data := <-dataChan

		jsonValue := make(map[string]bigquery.JsonValue)

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

		// Got a data point
		// TODO: Send data to Fluentd (or directly to BigQuery?)
		_, err = s.Tabledata.InsertAll(*bqproject, *bqdataset, *bqtable, &bigquery.TableDataInsertAllRequest{
			Kind: "bigquery#tableDataInsertAllRequest",
			Rows: []*bigquery.TableDataInsertAllRequestRows{
				&bigquery.TableDataInsertAllRequestRows{
					Json: jsonValue,
				},
			},
		}).Do()

		if err != nil {
			Error.Printf("Could not send data from %s to BigQuery: %v", data.Device.Name, err)
		} else {
			Debug.Printf("Data sent to BigQuery (%s): %s", data.Device.Name, data.Data)
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

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, bigquery.BigqueryScope)
	if err != nil {
		Error.Fatal("Could not create Bigquery client.", err)
	}
	service, err := bigquery.New(client)
	if err != nil {
		Error.Fatal("Could not create Bigquery client.", err)
	}

	// Start the thead to process data.
	go processData(service)

	// Start the api server
	go apiServer()

	// Start the monitor
	monitorDevices()
}
