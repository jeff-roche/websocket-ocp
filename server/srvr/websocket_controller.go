package srvr

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type IncommingMsg struct {
	Source  string
	Payload []byte
}

type WebsocketServerController struct {
	// Registered clients
	clients map[*ServerClient]bool

	// Inbound messages from the client
	broadcast chan IncommingMsg

	// register requests from clients
	register chan *ServerClient

	// Unregister request from clients
	unregister chan *ServerClient

	// Shutdown the controller
	shutdown chan os.Signal
}

func newWebsocketHub() *WebsocketServerController {
	return &WebsocketServerController{
		broadcast:  make(chan IncommingMsg),
		register:   make(chan *ServerClient),
		unregister: make(chan *ServerClient),
		clients:    make(map[*ServerClient]bool),
		shutdown:   nil,
	}
}

func (h *WebsocketServerController) run() {
	h.shutdown = make(chan os.Signal)

	signal.Notify(h.shutdown, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if client.id != message.Source {
					select {
					case client.send <- []byte(fmt.Sprintf("%s: %s", message.Source, string(message.Payload))):
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		case <-h.shutdown:
			log.Println("Killing the controller")
			return
		}
	}
}
