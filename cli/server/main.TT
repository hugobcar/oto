package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/hugobcar/oto/pkg/protobuf/app"
)

var port = "6885" string

func main() {
	fmt.Printf("Starting Oto server on port: %s", port)
	c := make(chan *app.LogsRequest)
	go func() {
		for {
			message := <-c
			ReadReceivedData(message)
		}
	}()
	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	for {
		if conn, err := listener.Accept(); err == nil {
			go handleProtoClient(conn, c)
		} else {
			continue
		}
	}
}

func ReadReceivedData(data *app.LogsRequest) {
	msgItems := data.GetPodName()
	fmt.Println("Receiving data…")
	for _, item := range msgItems {
		fmt.Println(item)
	}
}

func handleProtoClient(conn net.Conn, c chan *app.LogsRequest) {
	fmt.Println("Connected!")
	defer conn.Close()
	var buf bytes.Buffer
	_, err := io.Copy(&buf, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	pdata := new(app.LogsRequest)
	err = proto.Unmarshal(buf.Bytes(), pdata)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	c <- pdata
}
