package client

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func NewApp(conf *Conf, in io.Reader, out io.Writer) (*App, error) {
	return &App{
		conf: conf,
		in:   in,
		out:  out,
		hub:  NewHub(conf.User),
	}, nil
}

type App struct {
	conf   *Conf
	socket *websocket.Conn
	in     io.Reader
	out    io.Writer
	hub    *Hub
}

// connect to server
func (a *App) connect() error {
	url := a.conf.Host + "/ws"
	socket, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return errors.Wrap(err, `failed to connect to the server`)
	}
	a.socket = socket
	if err := a.hub.OnConnect(); err != nil {
		return errors.Wrap(err, `failed to send message on connect`)
	}
	return nil
}

func (a *App) disconnect() {
	select {
	case data, ok := <-a.hub.Done:
		if !ok {
			log.Printf(`read on closed Done channel`)
			return
		}
		a.socket.WriteMessage(websocket.TextMessage, data)
		a.socket.WriteMessage(websocket.CloseMessage, []byte{})
	}
}

// listen for incoming messages
func (a *App) listen() {
	log.Printf(`listen`)
	for {
		_, msg, err := a.socket.ReadMessage()
		if err != nil {
			log.Printf(`failed to read message from server: %s`, err)
			a.hub.OnDisconnect()
			return
		}
		a.hub.OnServerMessage(msg)
	}
}

// send message to server
func (a *App) send() {
	log.Printf(`send`)
	for {
		select {
		case data, ok := <-a.hub.Outbound:
			if !ok {
				a.hub.OnDisconnect()
				return
			}
			a.socket.WriteMessage(websocket.TextMessage, data)
		}
	}
}

// print incoming messages
func (a *App) print() {
	log.Printf(`print`)
	for {
		select {
		case data, ok := <-a.hub.Inbound:
			if !ok {
				a.hub.OnDisconnect()
				return
			}
			if _, err := a.out.Write([]byte(data)); err != nil {
				log.Printf(`failed to write message: %s`, err)
			}
		}
	}
}

// read user input
func (a *App) read() {
	log.Printf(`read`)
	scanner := bufio.NewScanner(a.in)
	for scanner.Scan() {
		msg := scanner.Bytes()
		if err := a.hub.OnUserMessage(msg); err != nil {
			log.Printf(`failed to read user message: %s`, err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf(`failed to read user messages: %s`, err)
	}
}

// Run client application
func (a *App) Run() error {
	// send messages to server
	go a.send()
	// print incoming messages
	go a.print()
	// read user input (or generate spam)
	go a.read()

	if err := a.connect(); err != nil {
		return errors.Wrap(err, `failed to connect`)
	}

	// stop on interrupt
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// listen server
	go a.listen()

	go func() {
		<-sigs
		a.hub.OnDisconnect()
	}()

	defer a.hub.Close()
	a.disconnect()
	return nil
}
