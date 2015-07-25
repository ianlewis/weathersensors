// Original Spark Community Thread: http://community.spark.io/t/how-to-access-the-core-via-we-browser/9711
// Code adapted from: http://arduino.cc/en/Tutorial/WebServer

/* Includes ------------------------------------------------------------------*/  
#include "application.h"

const String VERSION = String("1.0");
String deviceName = String("Unknown");

void log(String msg) {
    Serial.println(String("[") + String(Time.now()) + String("] ") + msg);
}

void deviceNameHandler(const char *topic, const char *data) {
    deviceName = String(data);
    log("Got device name: " + deviceName);
}

bool statusOk = false;

// How often we check the status in seconds.
const int STATUS_INTERVAL = 10;
const int NUM_LOOPS = 1000;
int loops = 0;
int nextStatus = 0;


// Update the status but only if we haven't
// printed the status recently.
bool checkStatus() {
    char ipAddress[15]; // holds the ip address
    int now;

    if (loops >= NUM_LOOPS) {
        now = Time.now();
        if (now >= nextStatus) {
            statusOk = WiFi.ready();
            if (statusOk) {
                // Once wifi is ready print the status and our IP address.
                IPAddress localIP = WiFi.localIP();
                sprintf(ipAddress, "%d.%d.%d.%d", localIP[0], localIP[1], localIP[2], localIP[3]);
            } else {
                sprintf(ipAddress, "<none>");
            }

            log("PING: DEVICE: " + deviceName + "; VERSION: " + VERSION + "; IP: " + ipAddress);

            nextStatus = now + STATUS_INTERVAL;
        }

        loops = 0;
    } else {
        loops += 1;   
    }

    return statusOk;
}


TCPServer server = TCPServer(5000);
TCPClient client;

const int REQ_MAX_BUF_SIZE = 1000;
char REQ_BUF[REQ_MAX_BUF_SIZE];
int REQ_BUF_SIZE = -1;

void serverMain() {
    // listen for incoming clients
    client = server.available();
    if (client) {
        // an http request ends with a blank line
        if (client.connected()) {
            log(String("Client connected."));
        }
        while (client.connected()) {
            // Check status while in this loop as well.
            checkStatus();

            if (client.available()) {
                REQ_BUF_SIZE += 1;
                REQ_BUF[REQ_BUF_SIZE] = client.read();
            }
            if (REQ_BUF[REQ_BUF_SIZE] == '\n') {
                // Replace the new line character with string termination.
                REQ_BUF[REQ_BUF_SIZE] = '\0';
                log(String(REQ_BUF));

                // Reset the request buffer.
                REQ_BUF_SIZE = -1;
            }
        }
        log(String("Client disconnected."));
    }
}

// Lists all sensors on this module.
void list() {

}

void get(String sensor) {

}

void setup() {
    // start listening for clients
    Serial.begin(115200);

    log("Getting device name...");
    Spark.subscribe("spark/device/name", deviceNameHandler);
    Spark.publish("spark/device/name");

    log("Starting server...");
    server.begin();
}



// The main loop that gets run forever.
void loop() {
    checkStatus();
    serverMain();
}
