// Original Spark Community Thread: http://community.spark.io/t/how-to-access-the-core-via-we-browser/9711
// Code adapted from: http://arduino.cc/en/Tutorial/WebServer

/* Includes ------------------------------------------------------------------*/  
#include "application.h"
#include "PietteTech_DHT.h"

const int READ_LED = D0;
const int DHTPIN = D1;
const int PORT = 5000;
const String VERSION = String("1.0");
String deviceName = String("");

void dht_wrapper(); // must be declared before the lib initialization

// Initialize the DHT22 sensor.
PietteTech_DHT dht(DHTPIN, DHT22, dht_wrapper);

// This wrapper is in charge of calling
// must be defined like this for the lib work
void dht_wrapper() {
    dht.isrCallback();
}

#define ONE_DAY_MILLIS (24 * 60 * 60 * 1000)
unsigned long lastSync = millis();

/*
 * Sync the device time with the spark server.
 */
void syncTime() {
    if (millis() - lastSync > ONE_DAY_MILLIS) {
        // Request time synchronization from the Particle Cloud
        Particle.syncTime();
        lastSync = millis();
    }
}

void log(String msg) {
    if (Serial.available()) {
        Serial.println(String("[") + String(Time.now()) + String("] ") + msg);
    }
}

void deviceNameHandler(const char *topic, const char *data) {
    deviceName = String(data);
    log("Got device name: " + deviceName);
}

void sendStatus(TCPClient client) {
    // Turn on the READ LED.
    digitalWrite(READ_LED, HIGH);

    float humidity = dht.readHumidity();
    float temp = dht.readTemperature();

    client.println(String("timestamp:") + String(Time.now()) + String("\ttemp:") + String(temp) + String("\thumidity:") + String(humidity));

    client.flush();

    // Delay so that the READ LED stays on
    // for a little longer.
    delay(100);
    digitalWrite(READ_LED, LOW);
}

TCPServer server = TCPServer(PORT);
TCPClient client;

void serverMain() {
    // listen for incoming clients
    client = server.available();
    if (client) {
        // an http request ends with a blank line
        if (client.connected()) {
            log(String("Client connected."));
            // Send LTSV to client.
            sendStatus(client);
        }
        client.stop();
        log(String("Client disconnected."));
    }
}


void setup() {
    // start listening for clients
    Serial.begin(115200);

    log("Getting device name...");
    Particle.subscribe("spark/device/name", deviceNameHandler);
    Particle.publish("spark/device/name");

    log("Starting server...");
    server.begin();

    log("Starting DHT22 sensor...");
    pinMode(READ_LED, OUTPUT);
}

// The main loop that gets run forever.
void loop() {
    // Sync with the Spark server if necessary.
    syncTime();

    // Run the main server loop.
    // NOTE: The server main function will not return while
    // a client has connected.
    serverMain();

    // Don't loop too fast.
    delay(100);
}
