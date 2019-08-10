package main

import (
	"github.com/oudwud/webtoon-service/pkg/config"

	"github.com/kelseyhightower/envconfig"
	"github.com/oudwud/webtoon-service/pkg/server"
	log "github.com/sirupsen/logrus"
)

const envPrefix = "WT"

func main() {
	var conf config.Config
	if err := envconfig.Process(envPrefix, &conf); err != nil {
		log.Fatal(err)
	}

	server.Run(&conf)
}
