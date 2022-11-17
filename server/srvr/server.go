package srvr

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//go:embed debug.html
var debugPage []byte

type WebsocketServer struct {
	addr     string
	hub      *WebsocketServerController
	listener chan error
	shutdown chan os.Signal
}

func NewWebsocketServer(addr string) *WebsocketServer {
	return &WebsocketServer{
		addr: addr,
		hub:  newWebsocketHub(),
	}
}

func (s WebsocketServer) Serve() {
	// Start the controller
	go s.hub.run()
	defer func() { close(s.hub.shutdown) }() // Shutdown the controller goroutine

	s.listener = make(chan error)
	s.shutdown = make(chan os.Signal)

	signal.Notify(s.shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	go s.listen()

	select {
	case err := <-s.listener:
		if err != nil {
			log.Fatal(err)
		}
	case <-s.shutdown:
		log.Println("Killing the server")
	}
}

func (s *WebsocketServer) listen() {
	// Setup the debug page
	http.HandleFunc("/debug", serveDebugPage)

	// Setup the client listeners
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(s.hub, w, r)
	})

	log.Printf("Listening for websocket connections on %s", s.addr)

	s.listener <- http.ListenAndServe(s.addr, nil)
}

func serveDebugPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Debug requested")

	if r.URL.Path != "/debug" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(200)
	_, err := w.Write(debugPage)
	if err != nil {
		log.Println("unable to serve the debug page:", err)
	}
}
