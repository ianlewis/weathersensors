/*
Firmware application for indoor_mod test sensor. The device includes
a DHT22 temperature/humidity sensor
*/

/* Includes ------------------------------------------------------------------*/  

#include "application.h"
#include "config.h"
#include "version.h"

// PietteTech DHT22 library.
// https://github.com/piettetech/PietteTech_DHT
#include "third_party/PietteTech_DHT/firmware/PietteTech_DHT.h"

// MQTT TLS library
// https://github.com/hirotakaster/MQTT-TLS
#include "third_party/MQTT-TLS/src/MQTT-TLS.h"

#define LETS_ENCRYPT_CA_PEM                                             \
"-----BEGIN CERTIFICATE----- \r\n"                                      \
"MIIFjTCCA3WgAwIBAgIRANOxciY0IzLc9AUoUSrsnGowDQYJKoZIhvcNAQELBQAw\r\n"  \
"TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh\r\n"  \
"cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMTYxMDA2MTU0MzU1\r\n"  \
"WhcNMjExMDA2MTU0MzU1WjBKMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg\r\n"  \
"RW5jcnlwdDEjMCEGA1UEAxMaTGV0J3MgRW5jcnlwdCBBdXRob3JpdHkgWDMwggEi\r\n"  \
"MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCc0wzwWuUuR7dyXTeDs2hjMOrX\r\n"  \
"NSYZJeG9vjXxcJIvt7hLQQWrqZ41CFjssSrEaIcLo+N15Obzp2JxunmBYB/XkZqf\r\n"  \
"89B4Z3HIaQ6Vkc/+5pnpYDxIzH7KTXcSJJ1HG1rrueweNwAcnKx7pwXqzkrrvUHl\r\n"  \
"Npi5y/1tPJZo3yMqQpAMhnRnyH+lmrhSYRQTP2XpgofL2/oOVvaGifOFP5eGr7Dc\r\n"  \
"Gu9rDZUWfcQroGWymQQ2dYBrrErzG5BJeC+ilk8qICUpBMZ0wNAxzY8xOJUWuqgz\r\n"  \
"uEPxsR/DMH+ieTETPS02+OP88jNquTkxxa/EjQ0dZBYzqvqEKbbUC8DYfcOTAgMB\r\n"  \
"AAGjggFnMIIBYzAOBgNVHQ8BAf8EBAMCAYYwEgYDVR0TAQH/BAgwBgEB/wIBADBU\r\n"  \
"BgNVHSAETTBLMAgGBmeBDAECATA/BgsrBgEEAYLfEwEBATAwMC4GCCsGAQUFBwIB\r\n"  \
"FiJodHRwOi8vY3BzLnJvb3QteDEubGV0c2VuY3J5cHQub3JnMB0GA1UdDgQWBBSo\r\n"  \
"SmpjBH3duubRObemRWXv86jsoTAzBgNVHR8ELDAqMCigJqAkhiJodHRwOi8vY3Js\r\n"  \
"LnJvb3QteDEubGV0c2VuY3J5cHQub3JnMHIGCCsGAQUFBwEBBGYwZDAwBggrBgEF\r\n"  \
"BQcwAYYkaHR0cDovL29jc3Aucm9vdC14MS5sZXRzZW5jcnlwdC5vcmcvMDAGCCsG\r\n"  \
"AQUFBzAChiRodHRwOi8vY2VydC5yb290LXgxLmxldHNlbmNyeXB0Lm9yZy8wHwYD\r\n"  \
"VR0jBBgwFoAUebRZ5nu25eQBc4AIiMgaWPbpm24wDQYJKoZIhvcNAQELBQADggIB\r\n"  \
"ABnPdSA0LTqmRf/Q1eaM2jLonG4bQdEnqOJQ8nCqxOeTRrToEKtwT++36gTSlBGx\r\n"  \
"A/5dut82jJQ2jxN8RI8L9QFXrWi4xXnA2EqA10yjHiR6H9cj6MFiOnb5In1eWsRM\r\n"  \
"UM2v3e9tNsCAgBukPHAg1lQh07rvFKm/Bz9BCjaxorALINUfZ9DD64j2igLIxle2\r\n"  \
"DPxW8dI/F2loHMjXZjqG8RkqZUdoxtID5+90FgsGIfkMpqgRS05f4zPbCEHqCXl1\r\n"  \
"eO5HyELTgcVlLXXQDgAWnRzut1hFJeczY1tjQQno6f6s+nMydLN26WuU4s3UYvOu\r\n"  \
"OsUxRlJu7TSRHqDC3lSE5XggVkzdaPkuKGQbGpny+01/47hfXXNB7HntWNZ6N2Vw\r\n"  \
"p7G6OfY+YQrZwIaQmhrIqJZuigsrbe3W+gdn5ykE9+Ky0VgVUsfxo52mwFYs1JKY\r\n"  \
"2PGDuWx8M6DlS6qQkvHaRUo0FMd8TsSlbF0/v965qGFKhSDeQoMpYnwcmQilRh/0\r\n"  \
"ayLThlHLN81gSkJjVrPI0Y8xCVPB4twb1PFUd2fPM3sA1tJ83sZ5v8vgFv2yofKR\r\n"  \
"PB0t6JzUA81mSqM3kxl5e+IZwhYAyO0OTg3/fs8HqGTNKd9BqoUwSRBzp06JMg5b\r\n"  \
"rUCGwbCUDI0mxadJ3Bz4WxR6fyNpBK2yAinWEsikxqEt\r\n"  \
"-----END CERTIFICATE----- "
const char letsencryptCaPem[] = LETS_ENCRYPT_CA_PEM;

const int LOOP_INTERVAL_SEC = 60;
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

double humidity = 0;
double temp = 0;
String deviceName = String("");
String location = String(CONFIG_LOCATION);
String localIP;
String deviceType = String("indoor_mod");
String version = String(gitversion);

struct State {
    bool wifiReady;
    bool particleConnected;
    bool mqttConnected;
    State() : wifiReady(false), particleConnected(false), mqttConnected(false) {}
    bool ready() {
        return wifiReady && particleConnected && mqttConnected; 
    }
};

State prevState, currentState;

// Initialized the client and set keepalive to 2 * loop interval
MQTT client(CONFIG_MQTT_HOST, CONFIG_MQTT_PORT, LOOP_INTERVAL_SEC*2, NULL);

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
    if (Serial.isConnected()) {
        Serial.println(String("[") + String(Time.now()) + String("] ") + msg);
    }
}

void deviceNameHandler(const char *topic, const char *data) {
    deviceName = String(data);
    log("Got device name: "+deviceName);
}

void setup() {
    // start listening for clients
    Serial.begin(115200);

    pinMode(READ_LED, OUTPUT);

    Particle.variable("deviceType", deviceType);
    Particle.variable("humidity", humidity);
    Particle.variable("location", location);
    Particle.variable("temperature", temp);
    Particle.variable("localIP", localIP);
    Particle.variable("version", version);

    // Get the device name
    Particle.subscribe("particle/device/name", deviceNameHandler);
    Particle.publish("particle/device/name");

    // Enable tls. set Root CA pem file.
    client.enableTls(letsencryptCaPem, sizeof(letsencryptCaPem));

    // Delay 10 seconds so we can connect for debugging
    waitFor(Serial.isConnected, 10000);

    log("setup finished");
}

// The main loop that gets run forever.
void loop() {
    int extra = 0;

    // Sync with the Spark server if necessary.
    syncTime();

    // Wait for the deviceName before looping
    if (deviceName == "") {
        log("Waiting for device name...");
        delay(100);
        return;
    }

    // Run the mqtt client loop
    client.loop();

    currentState = State();
    currentState.wifiReady = WiFi.ready();
    currentState.particleConnected = Particle.connected();
    currentState.mqttConnected = client.isConnected();

    if (!currentState.mqttConnected) {
        // Attempt to re-connect to the MQTT server
        log("Connecting to mqtt server");
        if (client.connect(deviceName, CONFIG_MQTT_USERNAME, CONFIG_MQTT_PASSWORD)) {
            log("Connected to mqtt server");
        } else {
            log("Failed to connect to mqtt server");
        }
        currentState.mqttConnected = client.isConnected();
    }

    // Log changes
    if (!prevState.wifiReady && currentState.wifiReady) {
        log("Connected to WiFi");
    }
    if (!prevState.particleConnected && currentState.particleConnected) {
        log("Connected to Particle");
    }
    prevState = currentState;

    // If not ready then blink
    if (!currentState.ready()) {
        digitalWrite(READ_LED, HIGH);
        delay(1000);
        digitalWrite(READ_LED, LOW);
        delay(1000);
        return;
    }

    localIP = String(WiFi.localIP());

    // Turn on the READ LED.
    digitalWrite(READ_LED, HIGH);

    float fHumidity = dht.readHumidity();
    float fTemp = dht.readTemperature();

    // Sometimes the DHT sensor returns -4 for some reason.
    // Retry up to 3 times.
    int retries = 0;
    while ((fTemp == -4 || fHumidity == -4) && retries < 3) {
        // The DHT library only reads the sensor every 2 seconds.
        // Wait over 2 seconds before attempting another read.
        delay(2001);
        extra += 2001;

        fHumidity = dht.readHumidity();
        fTemp = dht.readTemperature();

        retries += 1;
    }

    // Ignore values that are exactly -4.
    if (fTemp != -4 && fHumidity != -4) {
        temp = fTemp;
        humidity = fHumidity;

        String data = String("{\"location\":\""+location+"\",\"timestamp\":" + String(Time.now()) + ",\"temperature\":" + String(temp) +",\"humidity\":" + String(humidity) + "}");
        client.publish("home/" + location + "/climate", data);

        log(data);
    }

    // Delay so that the READ LED stays on
    // for a little longer.
    delay(100);
    extra += 100;
    digitalWrite(READ_LED, LOW);

    // Send data every interval
    delay(1000 * LOOP_INTERVAL_SEC - extra);
}
