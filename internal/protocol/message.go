package protocol

import (
	"errors"
	"fmt"
	"strings"
)

const (
	MESSAGE_AUTH = iota
	MESSAGE_CHAT
	MESSAGE_END

	DELIMITER  = "::"
	NEWLINE    = "\r\n"
	TOKEN_AUTH = "auth"
	TOKEN_END  = "end"
)

type Message struct {
	t           int
	commandCode string
	message     string
}

func (m *Message) Marshal() []byte {
	s := strings.Join([]string{m.commandCode, m.message}, DELIMITER) + NEWLINE
	return []byte(s)
}

// NewAuthMessage returns message that should be first in session
func NewAuthMessage(user string) (*Message, error) {
	user = strings.TrimSpace(user)
	if len(user) == 0 {
		return nil, errors.New("user name can not be empty")
	}
	return &Message{t: MESSAGE_AUTH, commandCode: TOKEN_AUTH, message: user}, nil
}

// NewChatMessage returns user message
func NewChatMessage(key, body string) (*Message, error) {
	return &Message{t: MESSAGE_CHAT, commandCode: key, message: body}, nil
}

// NewEndMessage returns message that should be latest in session
func NewEndMEssage() (*Message, error) {
	return &Message{t: MESSAGE_END, commandCode: TOKEN_END, message: ""}, nil
}

// ParseMessage recognizes message
func ParseMessage(data []byte) (*Message, error) {
	inp := string(data)
	inp = strings.TrimSpace(inp)
	s := strings.Split(inp, DELIMITER)
	switch len(s) {
	case 1:
		if s[0] == TOKEN_END {
			return NewEndMEssage()
		}
	case 2:
		if s[0] == TOKEN_AUTH {
			return NewAuthMessage(s[1])
		}
		return NewChatMessage(s[0], s[1])
	}
	return nil, fmt.Errorf(`failed to parse message: %s`, inp)
}

// IsUserMessage check type of message
func IsUserMessage(m *Message) bool {
	return m.t == MESSAGE_CHAT
}
