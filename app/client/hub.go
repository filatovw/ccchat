package client

import (
	"github.com/filatovw/ccchat/internal/protocol"
	"github.com/pkg/errors"
)

// NewHub creates Hub
func NewHub(user string) *Hub {
	return &Hub{
		Outbound: make(chan []byte),
		Inbound:  make(chan []byte),
		Done:     make(chan []byte),
		User:     user,
	}
}

// Hub is a middlware that manages channels
type Hub struct {
	Outbound chan []byte
	Inbound  chan []byte
	Done     chan []byte
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
		return errors.Wrap(err, `failed to process user message`)
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
	h.Done <- msg.Marshal()
	return nil
}
