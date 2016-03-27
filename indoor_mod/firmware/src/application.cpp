/*
Firmware application for indoor_mod test sensor. The device includes
a DHT22 temperature/humidity sensor
*/

/* Includes ------------------------------------------------------------------*/  

#include "application.h"

// PietteTech DHT22 library.
// https://github.com/piettetech/PietteTech_DHT
#include "PietteTech_DHT.h"

const int READ_LED = D0;
const int DHTPIN = D1;

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


double humidity = 0;
double temp = 0;
String localIP;
String deviceType = String("indoor_mod");

void setup() {
    // start listening for clients
    Serial.begin(115200);

    pinMode(READ_LED, OUTPUT);

    Particle.variable("deviceType", deviceType);
    Particle.variable("humidity", humidity);
    Particle.variable("temperature", temp);
    Particle.variable("localIP", localIP);

    // Delay 15 seconds so we can connect for debugging
    delay(15000);
}

// The main loop that gets run forever.
void loop() {
    // Sync with the Spark server if necessary.
    syncTime();

    // Turn on the READ LED.
    digitalWrite(READ_LED, HIGH);

    humidity = dht.readHumidity();
    temp = dht.readTemperature();
    localIP = String(WiFi.localIP());

    String data = String("timestamp:") + String(Time.now()) + String("\ttemp:") + String(temp) + String("\thumidity:") + String(humidity);
    Particle.publish("weatherdata", data);
    log(data);

    // Delay so that the READ LED stays on
    // for a little longer.
    delay(100);
    digitalWrite(READ_LED, LOW);

    // Send data every 1 min.
    delay(1000 * 60 - 100);
}
