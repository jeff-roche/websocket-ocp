package main

import (
	"flag"
	"websockserver/srvr"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	server := srvr.NewWebsocketServer(*addr)
	server.Serve()
}
