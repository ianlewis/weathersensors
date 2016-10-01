// poller.go implements the devicePoller that is used to poll the Particle
// Device API.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Device represents a Particle device.
type device struct {
	Id string `json:"id"`
	// The name may be blank for a time until
	// the reconciliation loop runs because device names
	// are not given in events.
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
}

// devicePoller encapsulates all the data needed for polling the Particle
// Device API
type devicePoller struct {
	accessToken string
	interval    int

	// The current list of known devices we are monitoring.
	devices []device

	// A list of device IDs the poller is monitoring. If empty all devices
	// returned from the API are monitored.
	deviceIds []string

	// A map from device ID to channel that cancels a device timeout.
	// TODO: cancelChan should probably be protected by a lock.
	cancelChan map[string]chan bool

	// The error channel. When devices time out they are sent to this channel.
	errorChan chan device

	done chan struct{}
	wg   *sync.WaitGroup
}

// newDevicePoller creates an new device poller. The poller can be stopped by
// closing the poller's done channel and waiting on the poller's waitgroup
func newDevicePoller(accessToken string, interval int, deviceIds []string) *devicePoller {
	return &devicePoller{
		accessToken: accessToken,
		interval:    interval,
		deviceIds:   deviceIds,
		cancelChan:  make(map[string]chan bool),
		errorChan:   make(chan device, 20),
		done:        make(chan struct{}),
		wg:          &sync.WaitGroup{},
	}
}

// monitorDevice returns true if the given device ID is being monitored
// otherwise it returns false.
func (p *devicePoller) monitorDevice(deviceId string) bool {
	if len(p.deviceIds) == 0 {
		return true
	}

	for _, id := range p.deviceIds {
		if id == deviceId {
			return true
		}
	}

	return false
}

// getDevices reads the device list from the Particle API
func (p *devicePoller) getDevices() ([]device, error) {
	var d []device

	req, err := http.NewRequest("GET", particleEndpoint, nil)
	if err != nil {
		return d, fmt.Errorf("Could not create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.accessToken)
	log.Printf("Getting devices...")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return d, fmt.Errorf("Error connecting to Particle API: %v", err)
	}

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(&d)
	if err != nil {
		return d, fmt.Errorf("Error decoding JSON data: %v", err)
	}

	return d, nil
}

// updateDevices compares the data received from the Particle Device API with
// what was received previously and starts a new timeout goroutine if a device has
// gone offline or stops a timeout if it has come back online.
func (p *devicePoller) updateDevices(newDevices []device) {
	// Check for new devices.
	for _, d := range newDevices {
		if !p.monitorDevice(d.Id) {
			continue
		}

		found := false

		for _, d2 := range p.devices {
			if d.Id == d2.Id {
				found = true
			}
		}
		if !found {
			if d.Connected {
				log.Printf("Device %s (%s) is online.", d.Name, d.Id)
				continue
			}

			if !d.Connected {
				log.Printf("Device %s (%s) is offline.", d.Name, d.Id)
				p.cancelChan[d.Id] = p.timeoutDevice(d)
			}
		}
	}

	for _, d := range p.devices {
		if !p.monitorDevice(d.Id) {
			continue
		}

		for _, d2 := range newDevices {

			if d.Id == d2.Id {
				// online -> offline
				if d.Connected && !d2.Connected {
					log.Printf("Device %s (%s) is offline.", d.Name, d.Id)
					p.cancelChan[d2.Id] = p.timeoutDevice(d2)
				}
				// offline -> online
				if !d.Connected && d2.Connected {
					log.Printf("Device %s (%s) is online.", d.Name, d.Id)

					// Stop the device timeout if there is one.
					if c, ok := p.cancelChan[d2.Id]; ok {
						c <- true
					}
				}
			}
		}
	}

	p.devices = newDevices
}

// timeoutDevice waits the deviceTimeout period and if the device doesn't
// come back online then it times out and an error is created in Stackdriver
// Error Reporting.
func (p *devicePoller) timeoutDevice(d device) chan bool {
	cancel := make(chan bool)

	go func() {
		p.wg.Add(1)
		select {
		case <-time.After(time.Duration(*deviceTimeout) * time.Second):
			log.Printf("Device timed out: %s (%s)", d.Name, d.Id)

			// Send error to the error channel.
			p.errorChan <- d

			delete(p.cancelChan, d.Id)
			p.wg.Done()
			return
		case <-cancel:
			delete(p.cancelChan, d.Id)
			p.wg.Done()
			return
		}
	}()

	return cancel
}

// poll polls the Particle Device API devices and checks their online
// status.
func (p *devicePoller) poll() {
	p.wg.Add(1)

	// Get the devices the first time.
	devices, err := p.getDevices()
	if err != nil {
		log.Printf("Error getting devices: %v", err)
	}
	if err == nil {
		p.updateDevices(devices)
	}

	// Poll the device API for device status.
	for {
		select {
		case <-time.After(time.Duration(p.interval) * time.Second):
			devices, err := p.getDevices()
			if err != nil {
				log.Printf("Error getting devices: %v", err)
				continue
			}
			p.updateDevices(devices)
		case <-p.done:
			// Cancel all timeout goroutines
			for _, c := range p.cancelChan {
				c <- true
			}
			p.wg.Done()
			log.Printf("Stopped polling loop.")
			return
		}
	}
}
