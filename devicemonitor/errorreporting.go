// errorreporting.go implements error reporting. If an error is sent to the
// error reporter via the error channel a report is generated on Stackdriver Error
// Reporting.

package main

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/errors"
)

func handleErrors(projectId string, deviceChan chan device, done chan struct{}, wg *sync.WaitGroup) {
	errorsClient, err := errors.NewClient(context.Background(), projectId, "devicemonitor", VERSION)
	if err != nil {
		log.Fatal("Could not create Stackdriver Error Reporting client: %v", err)
	}

	wg.Add(1)
	for {
		select {
		case d := <-deviceChan:
			handleError(errorsClient, d)
			continue
		case <-done:
			wg.Done()
			return
		}
	}
}

// handleError makes the request to the StackDriver Error Reporting API
func handleError(errorsClient *errors.Client, d device) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Sending report for %s (%s)", d.Name, d.Id)
	errorsClient.Reportf(ctx, nil, "Device is offline: %s (%s)", d.Name, d.Id)
}
