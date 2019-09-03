package main

import (
	"flag"
	"github.com/ngaut/log"

	"github.com/YangKeao/reserve/config"
	"github.com/YangKeao/reserve/server"

	"github.com/juju/errors"
)

func main() {
	var configPath = flag.String("config", "", "set the config path")
	flag.Parse()

	if *configPath == "" {
		log.Error("config path needs to be set")
		return
	}

	config, err := config.ReadAndParse(*configPath)
	if err != nil {
		log.Errorf("error while parsing config: %s", errors.ErrorStack(err))
		return
	}

	server, err := server.New(config)
	if err != nil {
		log.Errorf("error while creating server: %s", errors.ErrorStack(err))
	}

	server.ListenAndServe()
}
