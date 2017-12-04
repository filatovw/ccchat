package server

import (
	"log"

	"github.com/filatovw/ccchat/internal/protocol"
)

func NewHub() *Hub {
	return &Hub{
		clients:    NewClientsMap(),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

type Hub struct {
	clients    ClientsMap
	register   chan *Client
	unregister chan *Client
}

func (hub *Hub) addClient(c *Client) {
	hub.register <- c
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) send(msg *protocol.Message, client *Client) {
	client.outbound <- msg.Marshal()
}

func (hub *Hub) broadcast(msg *protocol.Message, ignore *Client) {
	for _, c := range hub.clients.data {
		if c.id != ignore.id {
			hub.send(msg, c)
		}
	}
}

func (hub *Hub) onConnect(c *Client) {
	hub.clients.Add(c)
	log.Println("client connected: ", c.socket.RemoteAddr())

}

func (hub *Hub) onDisconnect(c *Client) {
	hub.clients.Delete(c.id)
	c.close()
	log.Println("client disconnected: ", c.socket.RemoteAddr())
}

func (hub *Hub) onMessage(data []byte, c *Client) {
	msg, err := protocol.ParseMessage(data)
	if err != nil {
		log.Printf(`failed to parse message: %v`, msg)
	}
	log.Printf(`%s: %s`, c.id, string(msg.Marshal()))
	if protocol.IsAuthMessage(msg) {
		c.active = true
	} else if protocol.IsUserMessage(msg) && c.active {
		hub.broadcast(msg, c)
	} else if protocol.IsEndMessage(msg) && c.active {
		c.active = false
	} else {
		log.Printf(`unknown type of incoming message: %v`, msg)
	}
}
