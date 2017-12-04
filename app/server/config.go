package server

import (
	"flag"
	"log"
)

var serverConf = &ServerConf{}

func init() {
	flag.StringVar(&serverConf.Host, "host", "", "server host")
	flag.IntVar(&serverConf.Concurrency, "concurrency", 0, "pool size")
	flag.Parse()
}

type ServerConf struct {
	Host        string `json:"host"`
	Concurrency int    `json:"concurrency"`
}

func NewServerConf() *ServerConf {
	conf := &ServerConf{}
	if serverConf.Host == "" {
		conf.Host = "0.0.0.0:9000"
	}

	if serverConf.Concurrency <= 0 {
		conf.Concurrency = 100
	}
	log.Printf(`new server config: %+v`, conf)
	return conf
}
