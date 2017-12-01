package client

import (
	"fmt"

	"github.com/filatovw/ccchat/internal/protocol"
)

// NewHub creates Hub
func NewHub(user string) *Hub {
	return &Hub{
		Outbound: make(chan []byte),
		Inbound:  make(chan []byte),
		Done:     make(chan struct{}),
		User:     user,
	}
}

// Hub is a middlware that manages channels
type Hub struct {
	Outbound chan []byte
	Inbound  chan []byte
	Done     chan struct{}
	User     string
}

// Close closes all channels
func (h *Hub) Close() {
	close(h.Outbound)
	close(h.Inbound)
	close(h.Done)
}

// OnServerMessage send formatted messages
func (h Hub) OnServerMessage(d []byte) {
	h.Inbound <- d
}

// OnUserMessage send user from user to server
func (h Hub) OnUserMessage(d []byte) error {
	msg, err := protocol.ParseMessage(d)
	if err != nil {
		return fmt.Errorf(`failed to process user message: %s`, err)
	}
	h.Outbound <- msg.Marshal()
	return nil
}

// OnConnect sent auth message to server
func (h Hub) OnConnect() error {
	msg, err := protocol.NewAuthMessage(h.User)
	if err != nil {
		return err
	}
	h.Outbound <- msg.Marshal()
	return nil
}

// OnDisconnect send end message to server
func (h Hub) OnDisconnect() error {
	msg, err := protocol.NewEndMEssage()
	if err != nil {
		return err
	}
	h.Outbound <- msg.Marshal()
	return nil
}
