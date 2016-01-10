/*
Firmware application for outdoor_mod weather station.

The device includes a AM2315 temperature/humidity sensor which communicates
using i2c. i2c is implemented as a firmware library in Wire.h
*/

// AM2315 doesn't seem to be working. Disable for now.
#define ENABLE_AM2315 true
#define ENABLE_BMP180 true

/* Includes ------------------------------------------------------------------*/  

#include "application.h"

#if ENABLE_AM2315
// https://github.com/adafruit/Adafruit_AM2315
#include "Adafruit_AM2315.h"
#endif

#if ENABLE_BMP180
#include "Adafruit_BMP085_U.h"
#endif

#if ENABLE_AM2315
// Initialize the AM2315
Adafruit_AM2315 am2315;
#endif

#if ENABLE_BMP180
// Initialize the BMP180
Adafruit_BMP085_Unified bmp = Adafruit_BMP085_Unified(10085);
#endif

const int READ_LED = D2;

#define ONE_DAY_MILLIS (24 * 60 * 60 * 1000)
unsigned long lastSync = millis();

/*
 * Sync the device time with the spark server.
 */
void syncTime() {
    if (millis() - lastSync > ONE_DAY_MILLIS) {
        log("Syncing time with cloud...");
        // Request time synchronization from the Particle Cloud
        Particle.syncTime();
        lastSync = millis();
        log("Syncing time with cloud...");
    }
}

void log(String msg) {
    if (Serial.available()) {
        Serial.println(String("[") + String(Time.now()) + String("] ") + msg);
    }
}


float temp, humidity;

void setup() {
    Serial.begin(115200);

    // Delay 15 seconds so we can connect for debugging
    delay(15000);

#if ENABLE_AM2315
    log("Detecting AM2315...");
    if (!am2315.begin()) {
        log("Sensor not found, check wiring & pullups!");
        while(1);
    }
    log("AM2315 Detected.");
#endif

#if ENABLE_BMP180
    log("Detecting BMP180...");
    if (!bmp.begin()) {
        /* There was a problem detecting the BMP085 ... check your connections */
        Serial.print("Ooops, no BMP085 detected ... Check your wiring or I2C ADDR!");
        while(1);
    }
    log("BMP180 Detected.");
#endif

    pinMode(READ_LED, OUTPUT);
}

void loop() {
    // Sync with the Spark server if necessary.
    syncTime();

    String data = String("timestamp:") + String(Time.now());

    // Turn on the READ LED.
    digitalWrite(READ_LED, HIGH);

#if ENABLE_AM2315
    am2315.readTemperatureAndHumidity(temp, humidity);

    data += String("\ttemp:") + String(temp) + String("\thumidity:") + String(humidity);
#endif

#if ENABLE_BMP180
    /* Get a new sensor event */ 
    sensors_event_t event;
    bmp.getEvent(&event);
 
    /* Barometric pressure is measure in hPa */
    if (event.pressure) {
        data += String("\tpressure:") + String(event.pressure);
    }
#endif

    Particle.publish("weatherdata", data);
    log(data);

    // Delay so that the READ LED stays on
    // for a little longer.
    delay(100);
    digitalWrite(READ_LED, LOW);

    delay(1000 * 60 - 100);
}
