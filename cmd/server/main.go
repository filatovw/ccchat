package main

import (
	"log"
	"os"

	"github.com/filatovw/ccchat/app/server"
)

func main() {
	conf := server.NewConf()
	app := server.NewApp(conf)
	if err := app.Run(); err != nil {
		log.Fatalf(`fail: %s`, err)
	}
	os.Exit(0)
}
