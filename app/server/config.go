package server

import (
	"flag"
	"log"
)

var serverConf = &ServerConf{}

type ServerConf struct {
	Host   string `json:"host"`
	DBHost string `json:"dbhost"`
	DBPass string `json:"dbpass"`
}

func init() {
	flag.StringVar(&serverConf.Host, "host", "", "server host")
	flag.StringVar(&serverConf.DBHost, "dbhost", "", "db server host")
	flag.StringVar(&serverConf.DBPass, "dbpass", "", "db server pass")
	flag.Parse()
}

func NewServerConf() *ServerConf {
	conf := &ServerConf{}
	conf.Host = "0.0.0.0:9000"

	if serverConf.Host != "" {
		conf.Host = serverConf.Host
	}
	if serverConf.DBHost != "" {
		conf.DBHost = serverConf.DBHost
	}
	if serverConf.DBPass != "" {
		conf.DBPass = serverConf.DBPass
	}
	log.Printf(`new server config: %+v`, conf)
	return conf
}
