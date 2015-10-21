package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/donovanhide/eventsource"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/najeira/ltsv"
)

const VERSION = "0.6"

const PARTICLE_API_URL = "https://api.particle.io/v1/devices/events/weatherdata"

// TODO: Use FlagSet
var (
	fluentdHost       = flag.String("fluentd-host", "localhost", "The fluentd host.")
	fluentdPort       = flag.Int("fluentd-port", 24224, "The fluentd port.")
	fluentdRetryWait  = flag.Int("fluentd-retry", 500, "Amount of time is milliseconds to wait between retries.")
	debugLogging      = flag.Bool("debug", false, "Enable debug logging.")
	accessTokenPath   = flag.String("access-token-path", "", "The path to a file containing the Particle API access token.")
	particleRetryWait = flag.Int("particle-retry", 500, "Amount of time is milliseconds to wait between retries.")
	version           = flag.Bool("version", false, "Print the version and exit.")
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLogging() {
	debugOut := ioutil.Discard
	if *debugLogging {
		debugOut = os.Stdout
	}
	Debug = log.New(debugOut, "[DEBUG] ", log.Ldate|log.Ltime)
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	Warning = log.New(os.Stderr, "[WARNING] ", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
}

func getAccessToken() string {
	f, err := os.Open(*accessTokenPath)
	if err != nil {
		Error.Fatal("Could not open access token file: ", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		Error.Fatal("Could not open access token file: ", err)
	}
	return string(b)
}

// Data message from the Particle API
type Message struct {
	Id          string `json:"coreid"`
	Data        string `json:"data"`
	Ttl         string `json:"ttl"`
	PublishedAt string `json:"published_at"`
}

func addFloatValue(name string, jsonValue map[string]interface{}, data map[string]string) {
	if data[name] != "" {
		if val, err := strconv.ParseFloat(data[name], 64); err == nil {
			jsonValue[name] = val
		} else {
			Error.Printf("Error parsing %s data: %v", name, err)
		}
	}
}

// Continuously tries to connect to Fluentd.
func connectToFluentd() *fluent.Fluent {
	var err error
	var logger *fluent.Fluent

	// Continuously try to connect to Fluentd.
	for {
		Debug.Printf("Connecting to Fluentd (%s:%d)...", *fluentdHost, *fluentdPort)
		logger, err = fluent.New(fluent.Config{
			FluentHost: *fluentdHost,
			FluentPort: *fluentdPort,
			// Once we have a connection, the library will reconnect automatically
			// if the connection is lost. However, it panics if it fails to connect
			// more than MaxRetry times. To avoid panics crashing the server, retry
			// many times before panicking.
			MaxRetry:  240,
			RetryWait: *fluentdRetryWait,
		})
		if err != nil {
			Error.Printf("Could not connect to Fluentd: %v", err)
			time.Sleep(time.Duration(*fluentdRetryWait) * time.Millisecond)
		} else {
			return logger
		}
	}
}

func connectToParticle(accessToken string) *eventsource.Stream {
	// Continuously try to connect to the stream.
	for {
		req, err := http.NewRequest("GET", PARTICLE_API_URL, nil)
		if err != nil {
			Error.Fatal("Could not create request: %v", err)
			time.Sleep(time.Duration(*particleRetryWait) * time.Millisecond)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		Debug.Printf("Connecting to Particle API...")
		stream, err := eventsource.SubscribeWithRequest("", req)
		if err != nil {
			Error.Printf("Could not subscribe to Particle API stream: %v", err)
			time.Sleep(time.Duration(*particleRetryWait) * time.Millisecond)
		} else {
			return stream
		}
	}
}

// Proceses data in parallel
func processData(accessToken string) {
	var err error

	// Connect to Fluentd
	logger := connectToFluentd()

	// The stream object reconnects with exponential backoff.
	stream := connectToParticle(accessToken)

	// Now actually process events.
	for {
		// Block on the data channel.
		event := <-stream.Events

		// Unmarshall the JSON data from the Particle API.
		//////////////////////////////////////////////////////////////////
		var m Message
		jsonData := event.Data()
		// The particle API often sends newlines.
		// Perhaps as a keep-alive mechanism.
		if jsonData == "" {
			continue
		}
		err = json.Unmarshal([]byte(jsonData), &m)
		if err != nil {
			Error.Printf("Could not parse message data: %v", err)
			continue
		}

		// Read LTSV data from the device into map[string]string
		//////////////////////////////////////////////////////////////////
		reader := ltsv.NewReader(bytes.NewBufferString(m.Data))
		records, err := reader.ReadAll()
		if err != nil || len(records) != 1 {
			Error.Printf("Error reading LTSV data: %v", err)
			continue
		}

		data := records[0]

		// Put the data into jsonValue and send to BigQuery
		//////////////////////////////////////////////////////////////////
		jsonValue := make(map[string]interface{})

		jsonValue["deviceid"] = m.Id

		timestamp, err := strconv.ParseInt(data["timestamp"], 10, 64)
		if err != nil {
			Error.Printf("Error reading timestamp: %v", err)
			continue
		}

		jsonValue["timestamp"] = timestamp
		addFloatValue("temp", jsonValue, data)
		addFloatValue("humidity", jsonValue, data)
		addFloatValue("pressure", jsonValue, data)
		addFloatValue("windspeed", jsonValue, data)
		addFloatValue("winddirection", jsonValue, data)
		addFloatValue("rainfall", jsonValue, data)

		// Send data directly to Fluentd
		if err = logger.Post("aggre_mod.sensordata", jsonValue); err != nil {
			Error.Printf("Could not send data from %s to Fluentd: %v", m.Id, err)
		} else {
			Debug.Printf("Data processed (%s): %s", m.Id, data)
		}
	}
}

// A simple io.Writer wrapper around a logger so that we can use
// the logger as an io.Writer
type LogWriter struct{ *log.Logger }

func (w LogWriter) Write(b []byte) (int, error) {
	w.Printf("%s", b)
	return len(b), nil
}

func main() {
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	initLogging()

	// Get the API access token
	accessToken := getAccessToken()

	// Process data in the main thread.
	processData(accessToken)
}
