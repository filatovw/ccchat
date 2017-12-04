package server

import (
	"flag"
	"log"
)

var cliConf = &Conf{}

type Conf struct {
	Host   string `json:"host"`
	DBHost string `json:"dbhost"`
	DBPass string `json:"dbpass"`
}

func init() {
	flag.StringVar(&cliConf.Host, "host", "", "server host")
	flag.StringVar(&cliConf.DBHost, "dbhost", "", "db server host")
	flag.StringVar(&cliConf.DBPass, "dbpass", "", "db server pass")
	flag.Parse()
}

func NewConf() *Conf {
	conf := &Conf{}
	conf.Host = "0.0.0.0:9000"

	if cliConf.Host != "" {
		conf.Host = cliConf.Host
	}
	if cliConf.DBHost != "" {
		conf.DBHost = cliConf.DBHost
	}
	if cliConf.DBPass != "" {
		conf.DBPass = cliConf.DBPass
	}
	log.Printf(`new server config: %+v`, conf)
	return conf
}
