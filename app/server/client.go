package server

import (
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

// Client represents one real client
type Client struct {
	id       string
	name     string
	socket   *websocket.Conn
	hub      *Hub
	outbound chan []byte
}

func (client *Client) setName(name string) {
	client.name = name
}

func newClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		id:       uuid.NewV4().String(),
		name:     "anonymous",
		socket:   socket,
		hub:      hub,
		outbound: make(chan []byte),
	}
}

func (client *Client) read() {
	var (
		err  error
		data []byte
	)
	for {
		_, data, err = client.socket.ReadMessage()
		if err != nil {
			client.hub.onDisconnect(client)
			break
		}
		client.hub.onMessage(data, client)
	}
}

func (client *Client) write() {
	var (
		err  error
		data []byte
		ok   bool
	)
	for {
		select {
		case data, ok = <-client.outbound:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				client.socket.Close()
				return
			}
			if err = client.socket.WriteMessage(websocket.TextMessage, data); err != nil {
				client.hub.onDisconnect(client)
				break
			}
		}
	}
}

func (client Client) run() {
	go client.read()
	go client.write()
}

func (client Client) close() {
	close(client.outbound)
}
