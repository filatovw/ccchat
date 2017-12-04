package main

import (
	"io"
	"log"
	"os"

	"github.com/filatovw/ccchat/app/client"
	"github.com/filatovw/ccchat/internal/spammer"
)

func main() {
	conf, err := client.NewConf()
	if err != nil {
		log.Fatalf(`failed on configuration: %v`, err)
	}

	var inp io.Reader
	inp = os.Stdin
	if conf.GenConf != nil {
		inp = spammer.NewSpammer(
			conf.GenConf.Duration, conf.GenConf.UpperCase, conf.GenConf.Number)
	}

	app, err := client.NewApp(conf, inp, os.Stdout)
	if err != nil {
		log.Fatalf("failed on init:\n%+v\n", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("failed on run:\n%+v\n", err)
	}
	os.Exit(0)
}
