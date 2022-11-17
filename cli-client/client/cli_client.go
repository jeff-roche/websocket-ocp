package client

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type CliClient struct {
	url            string
	connection     *websocket.Conn
	interrupt      chan os.Signal
	socketReadDone chan interface{}
	cliReadDone    chan interface{}
	writer         io.Writer
	reader         io.Reader
	cmdPath        string
}

func NewCliClient(addr, path string) *CliClient {
	return &CliClient{
		url: "ws://" + addr + path,
	}
}

func (c *CliClient) Connect() {
	// Setup the interrupt handler
	c.interrupt = make(chan os.Signal)
	signal.Notify(c.interrupt, os.Interrupt)

	// Setup the websocket connection
	var err error
	c.connection, _, err = websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		log.Fatal("Error connecting to the websocket server:", err)
	} else {
		log.Println("Connected to the server at", c.url)
	}

	defer c.connection.Close()

	// Setup the cli io
	clireader, cliwriter, err := os.Pipe()
	if err != nil {
		log.Fatal("Unable to setup cli io pipe:", err)
	}

	c.reader, c.writer = clireader, cliwriter
	defer clireader.Close()
	defer cliwriter.Close()

	c.socketReadDone = make(chan interface{})

	// Listen for incomming messages
	go c.listenSock()

	// Listen for user input
	c.cliReadDone = make(chan interface{})
	go c.listenUser()

	// Primary event loop
	for {
		select {
		case <-c.socketReadDone:
			log.Println("Socket reading closed unexpectedly. Exiting...")
			err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error closing websocket connection:", err)
			}

			return

		case <-c.cliReadDone:
			log.Println("Stdin reading closed unexpectedly. Exiting...")
			err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error closing websocket connection:", err)
			}

			return

		case <-c.interrupt:
			log.Println("Interrupt received. Exiting...")

			err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error closing websocket connection:", err)
			}

			// Wait for the channel to close
			<-c.socketReadDone

			return
		}
	}

}

// Goroutine for listening to the websocket and forwarding to stdout
func (c *CliClient) listenSock() {
	for {
		_, msg, err := c.connection.ReadMessage()
		if err != nil {
			log.Println("Connection closed.")
			break
		}

		log.Println(string(msg))
	}

	close(c.socketReadDone)

}

// Goroutine for listening to usr input and forwarding to the websocket
func (c *CliClient) listenUser() {
	defer func() {
		close(c.cliReadDone)
	}()

	// Listen to standard in
	r := bufio.NewReader(os.Stdin)

	for {
		data, _, err := r.ReadLine()
		if err != nil {
			log.Println("Error reading from stdin:", err)
			break
		}

		err = c.connection.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Error writing message to websocket:", err)
			return
		}
	}
}
