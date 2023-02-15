package main

import (
	"go-protobuf-tcp/client"
	// "go-protobuf-tcp/server"
	"net"
)

func main() {
	c := client.NewClient(net.ParseIP("127.0.0.1"), 9527)
	err := c.Connect()
	if err != nil {
		panic(err)
	}
	c.Start()
	// server.NewServer(net.ParseIP("127.0.0.1"), 9527).Listen()
}
