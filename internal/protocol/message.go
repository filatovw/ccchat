package protocol

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DELIMITER  = "::"
	NEWLINE    = "\r\n"
	TOKEN_AUTH = "auth"
	TOKEN_END  = "end"
)

func marshal(command, message string) []byte {
	s := strings.Join([]string{command, message}, DELIMITER) + NEWLINE
	return []byte(s)
}

type Messager interface {
	Marshal() []byte
	MarshalServer(string) []byte
}

type AuthMessage struct {
	commandCode string
	message     string
}

func (m AuthMessage) Name() string {
	return m.message
}

func (m AuthMessage) Marshal() []byte {
	return marshal(m.commandCode, m.message)
}

func (m AuthMessage) MarshalServer(name string) []byte {
	return []byte(fmt.Sprintf(`[%s] %s`, name, m.commandCode))
}

// NewAuthMessage returns message that should be first in session
func NewAuthMessage(user string) (Messager, error) {
	user = strings.TrimSpace(user)
	if len(user) == 0 {
		return nil, errors.New("user name can not be empty")
	}
	return &AuthMessage{commandCode: TOKEN_AUTH, message: user}, nil
}

type UserMessage struct {
	commandCode string
	message     string
}

func (m UserMessage) Marshal() []byte {
	return marshal(m.commandCode, m.message)
}

func (m UserMessage) MarshalServer(name string) []byte {
	return []byte(fmt.Sprintf(`[%s] %s | %s`, name, m.commandCode, m.message))
}

// NewUserMessage returns user message
func NewUserMessage(key, body string) (Messager, error) {
	return &UserMessage{commandCode: key, message: body}, nil
}

type EndMessage struct {
	commandCode string
}

func (m EndMessage) Marshal() []byte {
	return []byte(m.commandCode)
}

func (m EndMessage) MarshalServer(name string) []byte {
	return []byte(fmt.Sprintf(`[%s] %s`, name, m.commandCode))
}

// NewEndMessage returns message that should be latest in session
func NewEndMEssage() (Messager, error) {
	return &EndMessage{TOKEN_END}, nil
}

// ParseMessage recognizes message
func ParseMessage(data []byte) (Messager, error) {
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
		return NewUserMessage(s[0], s[1])
	}
	return nil, fmt.Errorf(`failed to parse message: %s`, inp)
}
