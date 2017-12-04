package client

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

type ClientConf struct {
	User    string  `json:"user"`
	Host    string  `json:"host"`
	GenConf GenConf `json:"gen,omitempty"`
}

type GenConf struct {
	Number    uint          `json:"number,omitempty"`
	UpperCase bool          `json:"upper_case,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
}

func (g *GenConf) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	cnf := &struct {
		Number    uint   `json:"number,omitempty"`
		UpperCase bool   `json:"upper_case,omitempty"`
		Duration  string `json:"duration,omitempty"`
	}{}
	err := json.Unmarshal(data, cnf)
	if err != nil {
		return err
	}

	ts, err := time.ParseDuration(cnf.Duration)
	if err != nil {
		return errors.Wrapf(err, `wrong duration format: "%s"`, cnf.Duration)
	}

	g.Number = cnf.Number
	g.UpperCase = cnf.UpperCase
	g.Duration = ts
	return nil
}

// NewConf creates new client configuration
func NewConf(conf, user, host string, genNumber uint, genUpperCase bool, genDuration string) (*ClientConf, error) {
	c := &ClientConf{GenConf: GenConf{}}
	if conf != "" {
		data, err := ioutil.ReadFile(conf)
		if err != nil {
			return nil, errors.Wrap(err, `failed to load config`)
		}
		if err := json.Unmarshal(data, c); err != nil {
			return nil, errors.Wrap(err, `failed to parse config`)
		}
	}
	if user != "" {
		c.User = user
	}
	if host != "" {
		c.Host = host
	}
	if genDuration != "" {
		duration, err := time.ParseDuration(genDuration)
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse duration`)
		}
		c.GenConf.Duration = duration
	}

	if genNumber != 0 {
		c.GenConf.Number = genNumber
	}
	c.GenConf.UpperCase = genUpperCase
	return c, nil
}
