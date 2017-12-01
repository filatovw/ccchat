package main

import (
	"flag"
	"log"
	"os"

	"github.com/filatovw/ccchat/app/client"
)

type Params struct {
	Conf         string // path to config
	User         string
	Host         string
	GenNumber    uint
	GenUpperCase bool
	GenDuration  string
}

var params = &Params{}

func init() {
	flag.StringVar(&params.Conf, "conf", "./client.json", "path to config (default: ./client.json)")
	flag.StringVar(&params.User, "user", "", "user name")
	flag.StringVar(&params.Host, "host", "0.0.0.0:9000", "server address [host]:[port]")
	flag.UintVar(&params.GenNumber, "gen.number", 0, "number of autogenerated messages")
	flag.BoolVar(&params.GenUpperCase, "gen.uppercase", false, "generated messages in upper case")
	flag.StringVar(&params.GenDuration, "gen.duration", "", "generation session duration")
	flag.Parse()
}

func main() {
	conf, err := client.NewConf(params.Conf, params.User, params.Host, params.GenNumber, params.GenUpperCase, params.GenDuration)
	if err != nil {
		log.Fatalf(`failed on configuration: %s`, err)
	}

	app := client.NewApp(conf, os.Stdin, os.Stdout)
	if err := app.Run(); err != nil {
		log.Fatalf(`fail: %s`, err)
	}
	os.Exit(0)
}
