package spammer

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpammerNumber(t *testing.T) {
	spammer := NewSpammer(time.Second, false, 1)
	p := make([]byte, 10, 10)
	size, err := spammer.Read(p)
	assert.NoError(t, err)
	assert.True(t, size > 0)
	size, err = spammer.Read(p)
	assert.Error(t, err)
	assert.True(t, size == 0)
}

func TestSpammerUppercase(t *testing.T) {
	hasAllUpper := func(b []byte) bool {
		s := string(b)
		for _, l := range s {
			stl := string(l)
			if stl != strings.ToUpper(stl) {
				return false
			}
		}
		return true
	}

	spammer := NewSpammer(time.Second, true, 1)
	p := make([]byte, 10, 10)
	size, err := spammer.Read(p)
	assert.NoError(t, err)
	assert.True(t, size > 0)

	assert.True(t, hasAllUpper(p))

	size, err = spammer.Read(p)
	assert.Error(t, err)
	assert.True(t, size == 0)
}

func TestSpammerTimeout(t *testing.T) {
	spammer := NewSpammer(time.Nanosecond, false, 1)
	p := make([]byte, 10, 10)
	size, err := spammer.Read(p)
	assert.NoError(t, err)
	assert.True(t, size > 0)

	time.Sleep(time.Nanosecond * 2)

	size, err = spammer.Read(p)
	assert.Error(t, err)
	assert.True(t, size == 0)
}
