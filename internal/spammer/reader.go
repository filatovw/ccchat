package spammer

import (
	"bytes"
	"io"
	"time"

	"github.com/Pallinder/go-randomdata"
)

type Spammer struct {
	duration  *time.Timer
	uppercase bool
	num       uint

	counter uint
}

func NewSpammer(d time.Duration, up bool, num uint) *Spammer {
	s := &Spammer{uppercase: up, num: num}
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
		g := s.gen()
		if s.uppercase {
			g = bytes.ToUpper(p)
		}
		copy(p, g)
		s.counter++
	}
	return len(p), nil
}
