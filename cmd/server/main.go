package main

import (
	"log"
	"os"

	"github.com/filatovw/ccchat/app/server"
)

func main() {
	// create config from CLI arguments
	conf, err := server.NewConf()
	if err != nil {
		log.Fatalf(`failed on configuration: %v`, err)
	}

	// create server app
	app := server.NewApp(conf)
	// start app
	if err := app.Run(); err != nil {
		log.Fatalf(`fail: %s`, err)
	}
	os.Exit(0)
}
