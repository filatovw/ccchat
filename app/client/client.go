package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

func NewApp(conf *ClientConf, in io.Reader, out io.Writer) *App {
	return &App{
		conf: conf,
		in:   in,
		out:  out,
		hub:  NewHub(conf.User),
	}
}

type App struct {
	conf   *ClientConf
	socket *websocket.Conn
	in     io.Reader
	out    io.Writer
	hub    *Hub
}

// connect to server
func (a *App) connect() error {
	socket, _, err := websocket.DefaultDialer.Dial(a.conf.Host, nil)
	if err != nil {
		return fmt.Errorf(`failed to connect to the server: %s`, err)
	}
	a.socket = socket
	if err := a.hub.OnConnect(); err != nil {
		return fmt.Errorf(`failed to connect to the server: %s`, err)
	}
	return nil
}

// listen for incoming messages
func (a *App) listen() {
	for {
		_, msg, err := a.socket.ReadMessage()
		if err != nil {
			log.Printf(`failed to read message from server: %s`, err)
			return
		}
		a.hub.OnServerMessage(msg)
	}
}

// send message to server
func (a *App) send() {
	for {
		select {
		case data, ok := <-a.hub.Outbound:
			if !ok {
				a.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			a.socket.WriteMessage(websocket.TextMessage, data)
		case <-a.hub.Done:
			return
		}
	}
}

// print incoming messages
func (a *App) print() {
	for {
		select {
		case data, ok := <-a.hub.Inbound:
			if !ok {
				return
			}
			_, err := a.out.Write(data)
			log.Printf(`failed to write message: %s`, err)
		case <-a.hub.Done:
			return
		}
	}
}

// read user input
func (a *App) read() {
	scanner := bufio.NewScanner(a.in)
	for scanner.Scan() {
		a.out.Write([]byte("[key]::[message] >"))
		if err := a.hub.OnUserMessage(scanner.Bytes()); err != nil {
			log.Printf(`failed to read user message: %s`, err)
			a.out.Write([]byte("failed to read message\r\n"))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf(`failed to read user messages: %s`, err)
	}
}

// Run client application
func (a *App) Run() error {
	if err := a.connect(); err != nil {
		return fmt.Errorf(`failed to connect: %s`, err)
	}
	defer a.hub.Close()
	// stop on interrupt
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		a.hub.Done <- struct{}{}
	}()

	// listen server
	go a.listen()
	// send messages to server
	go a.send()
	// print incoming messages
	go a.print()
	// read user input (or generate spam)
	go a.read()

	<-a.hub.Done
	return nil
}
