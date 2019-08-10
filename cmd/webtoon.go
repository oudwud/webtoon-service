package main

import (
	"github.com/oudwud/webtoon-service/pkg/config"
	"github.com/oudwud/webtoon-service/pkg/server"
)

func main() {
	conf := config.New()
	server.Run(conf)
}
