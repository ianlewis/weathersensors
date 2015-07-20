// Original Spark Community Thread: http://community.spark.io/t/how-to-access-the-core-via-we-browser/9711
// Code adapted from: http://arduino.cc/en/Tutorial/WebServer

/* Includes ------------------------------------------------------------------*/  
TCPServer server = TCPServer(80);
TCPClient client;

void setup()
{
    // start listening for clients
    server.begin();
}

void loop() {
    // listen for incoming clients
    client = server.available();
    if (client) {
        // an http request ends with a blank line
        boolean currentLineIsBlank = true;
        while (client.connected()) {
            if (client.available()) {
                char c = client.read();
                // if you've gotten to the end of the line (received a newline
                // character) and the line is blank, the http request has ended,
                // so you can send a reply
                if (c == '\n' && currentLineIsBlank) {
                    // send a standard http response header
                    client.println("HTTP/1.1 200 OK");
                    client.println("Content-Type: application/json");
                    client.println("Connection: close");  // the connection will be closed after completion of the response
                    client.println();
                    client.print("{\"time\": \"");
                    client.print(Time.timeStr());
                    client.println("\"}");
                    client.println();
                    delay(5);
                    break;
                }
                if (c == '\n') {
                    // you're starting a new line
                    currentLineIsBlank = true;
                }
                else if (c != '\r') {
                    // you've gotten a character on the current line
                    currentLineIsBlank = false;
                }
            }
        }
        // give the web browser time to receive the data
        delay(1);
        // close the connection:
        client.stop();
    }
}
