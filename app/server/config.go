package server

import (
	"flag"
	"fmt"
	"log"
)

var cliConf = &Conf{}

type Conf struct {
	Host   string `json:"host"`
	DBHost string `json:"dbhost"`
	DBPass string `json:"dbpass"`
	DBUser string `json:"dbuser"`
	DBName string `json:"dbname"`
}

func (c Conf) String() string {
	return fmt.Sprintf("%T\n host: %s\n db.host: %s\n db.pass: %s\n db.user: %s\n db.name: %s",
		c, c.Host, c.DBHost, c.DBPass, c.DBUser, c.DBName)
}

func init() {
	flag.StringVar(&cliConf.Host, "host", "", "server host")
	flag.StringVar(&cliConf.DBHost, "dbhost", "", "db host")
	flag.StringVar(&cliConf.DBPass, "dbpass", "", "db pass")
	flag.StringVar(&cliConf.DBUser, "dbuser", "", "db user")
	flag.StringVar(&cliConf.DBName, "dbname", "", "db name")
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
	if cliConf.DBUser != "" {
		conf.DBUser = cliConf.DBUser
	}
	if cliConf.DBName != "" {
		conf.DBName = cliConf.DBName
	}
	log.Printf(`new server config: %s`, conf)
	return conf
}
