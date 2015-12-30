/*
Firmware application for outdoor_mod weather station.

The device includes a AM2315 temperature/humidity sensor which communicates
using i2c. i2c is implemented as a firmware library in Wire.h
*/

/* Includes ------------------------------------------------------------------*/  

#include "application.h"
// #include "Adafruit_AM2315.h"
#include "Adafruit_BMP085_U.h"

// Initialize the AM2315
// Adafruit_AM2315 am2315;

// Initialize the BMP180
Adafruit_BMP085_Unified bmp = Adafruit_BMP085_Unified(10085);


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

void setup() {
    Serial.begin(115200);

    // Delay 15 seconds so we can connect for debugging
    delay(15000);

    // while (!am2315.begin()) {
    //     log("Sensor not found, check wiring & pullups!");
    //     delay(1000);
    // }

    Serial.println("Detecting BMP180...");
    if (!bmp.begin()) {
        /* There was a problem detecting the BMP085 ... check your connections */
        Serial.print("Ooops, no BMP085 detected ... Check your wiring or I2C ADDR!");
        while(1);
    }
}

// char str[8];
// unsigned long result = 0;
// unsigned int temp, humi = 0;
// unsigned int amID = 0;

void loop() {
    // log("Hum: "); Serial.println(am2315.readHumidity());
    // log("Temp: "); Serial.println(am2315.readTemperature());

    /* Get a new sensor event */ 
    sensors_event_t event;
    bmp.getEvent(&event);
 
    /* Display the results (barometric pressure is measure in hPa) */
    if (event.pressure) {
        /* Display atmospheric pressure in hPa */
        Serial.print("Pressure: ");
        Serial.print(event.pressure);
        Serial.println(" hPa");

        Serial.print("Altitude: ");
        Serial.print(event.pressure);
        Serial.println(" hPa");
    } else {
        Serial.println("Sensor error");
    }

    delay(1000);
}
