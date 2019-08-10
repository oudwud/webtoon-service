package config

type Config struct {
	Debug bool `default:"false" envconfig:"DEBUG"`
	Port  int  `envconfig:"PORT"`
}
