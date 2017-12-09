package server

import (
	"log"

	"github.com/filatovw/ccchat/app/server/model"
	"github.com/filatovw/ccchat/internal/protocol"
	"github.com/jinzhu/gorm"
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
	db         *gorm.DB
}

func (hub *Hub) setDB(db *gorm.DB) {
	hub.db = db
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

func (hub *Hub) send(msg []byte, c *Client) {
	c.outbound <- msg
}

func (hub *Hub) broadcast(msg []byte, ignore *Client) {
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
		return
	}

	if m, ok := msg.(*protocol.AuthMessage); ok {
		c.name = m.Name()
		if _, err := model.GetOrCreateUser(hub.db, c.name); err != nil {
			log.Printf(`failed to get or create user %s`, err)
			return
		}

		hub.clients.Add(c)
	} else if _, ok := msg.(*protocol.EndMessage); ok {
		hub.unregister <- c
	}

	log.Printf("%s", msg.MarshalServer(c.name))
	hub.broadcast(msg.MarshalServer(c.name), c)

	user, err := model.GetOrCreateUser(hub.db, c.name)
	if err != nil {
		log.Printf(`failed to get or create user %s`, err)
		return
	}
	err = model.AddMessage(hub.db, user, string(msg.Marshal()))
	if err != nil {
		log.Printf(`failed to save message %s`, err)
		return
	}
}
