package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"io/ioutil"
)

// CLIConf place for arguments from command line
type CLIConf struct {
	Conf   string
	Host   string
	DBHost string
	DBPass string
	DBUser string
	DBName string
}

var cliConf = &CLIConf{}

// Conf storage for parameters of server
type Conf struct {
	Host string `json:"host"`
	DB   DBConf `json:"db"`
}

// DBConf is a storage for database parameters
type DBConf struct {
	Host string `json:"host"`
	Pass string `json:"pass"`
	User string `json:"user"`
	Name string `json:"name"`
}

// String implementation of Stringer interface
func (c Conf) String() string {
	return fmt.Sprintf("%T\n=====\n host: %s\n db.host: %s\n db.pass: %s\n db.user: %s\n db.name: %s\n=====",
		c, c.Host, c.DB.Host, c.DB.Pass, c.DB.User, c.DB.Name)
}

func init() {
	flag.StringVar(&cliConf.Conf, "conf", "", "path to config (./server.json)")
	flag.StringVar(&cliConf.Host, "host", "", "server host")
	flag.StringVar(&cliConf.DBHost, "db.host", "", "db host")
	flag.StringVar(&cliConf.DBPass, "db.pass", "", "db pass")
	flag.StringVar(&cliConf.DBUser, "db.user", "", "db user")
	flag.StringVar(&cliConf.DBName, "db.name", "", "db name")
	flag.Parse()
}

// NewConf creates server configuration
func NewConf() (*Conf, error) {
	c := &Conf{DB: DBConf{}}

	if cliConf.Conf != "" {
		data, err := ioutil.ReadFile(cliConf.Conf)
		if err != nil {
			return nil, errors.Wrap(err, `failed to load config`)
		}
		if err := json.Unmarshal(data, c); err != nil {
			return nil, errors.Wrap(err, `failed to parse config`)
		}
	}

	if cliConf.Host != "" {
		c.Host = cliConf.Host
	}
	if cliConf.DBHost != "" {
		c.DB.Host = cliConf.DBHost
	}
	if cliConf.DBPass != "" {
		c.DB.Pass = cliConf.DBPass
	}
	if cliConf.DBUser != "" {
		c.DB.User = cliConf.DBUser
	}
	if cliConf.DBName != "" {
		c.DB.Name = cliConf.DBName
	}
	log.Printf(`new server config: %s`, c)
	return c, nil
}
