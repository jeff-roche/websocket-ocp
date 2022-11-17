package main

import (
	"flag"
	"websockclient/client"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var path = flag.String("path", "/ws", "the address path for the websocket connection")

func main() {
	flag.Parse()

	client := client.NewCliClient(*addr, *path)

	client.Connect()
}
