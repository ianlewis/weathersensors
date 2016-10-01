// Command devicemonitor is a monitor that checks if particle devices are
// online and writes an error to Stackdriver Error Reporting.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

//go:generate go run scripts/gen.go

const (
	particleEndpoint = "https://api.particle.io/v1/devices"
)

var (
	addr = flag.String("host", stringDefaults(":8080", os.Getenv("ADDRESS")), "The web server address for health checks.")

	projectId = flag.String("project", os.Getenv("GCP_PROJECT"), "The Google Cloud Platform project ID for the Error Reporting API.")

	deviceTimeout = flag.Int("device-timeout", intDefaults(300, os.Getenv("DEVICE_TIMEOUT")), "The time that a device can be offline before an error is produced.")

	accessTokenPath = flag.String("access-token", os.Getenv("ACCESS_TOKEN_PATH"), "The path to a file containing the Particle API access token.")

	deviceListPath = flag.String("device-list", os.Getenv("DEVICE_LIST_PATH"), "The path to a text file of device IDs (one per line) to monitor. If not specified, all devices are monitored.")

	pollInterval = flag.Int("poll-interval", intDefaults(30, os.Getenv("POLL_INTERVAL")), "API polling interval in seconds.")
	version      = flag.Bool("version", false, "Print the version and exit.")
)

// stringDefaults takes a default value and a list of string values and returns the first
// non-empty value. If all values are empty or there are no values present
// the default string value is returned.
func stringDefaults(def string, val ...string) string {
	for i := range val {
		if val[i] != "" {
			return val[i]
		}
	}
	return def
}

// intDefaults takes a default int value and a list of string values and returns the first
// non-empty value that can be converted to an integer. If all values are empty
// or there are no values present the default int value is returned.
func intDefaults(def int, val ...string) int {
	for i := range val {
		if val[i] != "" {
			intVal, err := strconv.ParseInt(val[i], 10, 32)
			if err == nil {
				return int(intVal)
			}
		}
	}
	return def
}

// readAccessToken reads the access token from the path given on the command
// line.
func readAccessToken() string {
	f, err := os.Open(*accessTokenPath)
	if err != nil {
		log.Fatal("Could not open access token file: ", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("Could not open access token file: ", err)
	}
	return strings.Trim(string(b), " \t\n")
}

// readDeviceIds reads the text file of device IDs to monitor and updates the
// deviceIds global.
func readDeviceIds() []string {
	deviceIds := []string{}

	if *deviceListPath == "" {
		// No path was specified. Monitor all devices.
		return deviceIds
	}
	f, err := os.Open(*deviceListPath)
	if err != nil {
		log.Fatal("Error reading device ID list: ", err)
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		deviceIds = append(deviceIds, strings.Trim(s.Text(), " \t\n"))
	}

	if err := s.Err(); err != nil {
		log.Fatal("Error reading device ID list: ", err)
	}

	return deviceIds
}

// Returns the health status of the app.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

// Prints the server version
func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, VERSION)
}

func main() {
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	if *deviceTimeout < *pollInterval {
		log.Fatal("Device timeout cannot be less than the poll interval.")
	}

	// Get the API access token
	accessToken := readAccessToken()

	// Read the device ID list
	deviceIds := readDeviceIds()

	// Poll the device API for device status.
	poller := newDevicePoller(accessToken, *pollInterval, deviceIds)
	go poller.poll()

	go handleErrors(*projectId, poller.errorChan, poller.done, poller.wg)

	// Set up the web server for health checks.
	go func() {
		http.HandleFunc("/_status/healthz", healthHandler)
		http.HandleFunc("/_status/version", versionHandler)

		log.Printf("Listening on %s...", *addr)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(poller.done)
			poller.wg.Wait()
			log.Printf("Done.")
			os.Exit(0)
		}
	}
}
