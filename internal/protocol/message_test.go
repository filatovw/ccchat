package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageParse(t *testing.T) {
	check := func(input []byte, expected Messager) {
		m, err := ParseMessage(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, m)
	}

	msg, _ := NewUserMessage("aaa", "bbb")
	check([]byte("aaa::bbb"), msg)
	msg, _ = NewUserMessage("", "bbb")
	check([]byte("::bbb"), msg)
	msg, _ = NewUserMessage("aaa", "")
	check([]byte("aaa::"), msg)
	msg, _ = NewAuthMessage("user")
	check([]byte("auth::user"), msg)
	msg, _ = NewEndMEssage()
	check([]byte("end"), msg)
}

func TestMessageMarshal(t *testing.T) {
	check := func(input Messager, expected []byte) {
		assert.Equal(t, input.Marshal(), expected)
	}

	msg, _ := NewUserMessage("aaa", "bbb")
	check(msg, []byte("aaa::bbb\r\n"))
	msg, _ = NewUserMessage("aaa", "")
	check(msg, []byte("aaa::\r\n"))
	msg, _ = NewUserMessage("", "bbb")
	check(msg, []byte("::bbb\r\n"))
	msg, _ = NewUserMessage("", "")
	check(msg, []byte("::\r\n"))
	msg, _ = NewAuthMessage("user")
	check(msg, []byte("auth::user\r\n"))
}
