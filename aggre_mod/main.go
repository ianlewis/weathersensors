package main

import (
    "fmt"
    "net"
    "flag"
    "bufio"
    "strconv"

    "bitbucket.org/IanLewis/homesensorsproject/aggre_mod/internal/github.com/najeira/ltsv"
)

// TODO: Service discovery
var (
    addr = flag.String("addr", "127.0.0.1", "The host to connect to.")
    port = flag.Int("port", 5000, "The port number to connect on.")
)

func main() {
    flag.Parse()

    tcpAddr, err := net.ResolveTCPAddr("tcp", *addr + ":" + strconv.Itoa(*port))
    if err != nil {
        panic(err)
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    reader := ltsv.NewReader(bufio.NewReader(conn))

    for {
        if data, err := reader.Read(); err == nil {
            for k,v := range data {
                fmt.Print(k + ": " + v + "\t")
            }
            fmt.Println("")
        }
    }
}
