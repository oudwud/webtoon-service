package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Debug bool `default:"false" envconfig:"DEBUG"`
	Port  int  `required:"true" envconfig:"PORT"`
}

const envPrefix = "WT"

func New() *Config {
	var conf Config
	if err := envconfig.Process(envPrefix, &conf); err != nil {
		log.Fatal(err)
	}

	log.Info("DEBUG: ", conf.Debug)
	log.Info("PORT: ", conf.Port)

	if conf.Debug {
		log.SetLevel(log.DebugLevel)
	}

	return &conf
}
