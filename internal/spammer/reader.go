package spammer

import (
	"bytes"
	"io"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/filatovw/ccchat/internal/protocol"
)

type Spammer struct {
	duration  *time.Timer
	uppercase bool
	num       uint
	delay     time.Duration

	counter uint
}

func NewSpammer(d time.Duration, up bool, num uint, delay time.Duration) *Spammer {
	s := &Spammer{uppercase: up, num: num, delay: delay}
	if int(d) > 0 {
		s.duration = time.NewTimer(d)
	}
	return s
}

func (s Spammer) gen() []byte {
	return []byte(randomdata.SillyName())
}

func (s *Spammer) Read(p []byte) (int, error) {
	if s.counter >= s.num {
		return 0, io.EOF
	}
	select {
	case <-s.duration.C:
		return 0, io.EOF
	default:
		msg, _ := protocol.NewUserMessage(randomdata.Noun(), randomdata.Adjective())
		g := msg.Marshal()
		if s.uppercase {
			g = bytes.ToUpper(p)
		}
		copy(p, g)
		s.counter++
		time.Sleep(s.delay)
		return len(g), nil
	}
}
