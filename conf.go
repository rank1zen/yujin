package main

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Debug              bool
	ServerPort         int    `split_words:"true"`
	PostgresConnString string `split_words:"true" required:"true"`
}

var conf = Config{
	Debug:      true,
	ServerPort: 1323,
}

func LoadConfig() (*Config, error) {
	if err := envconfig.Process("YUJIN", &conf); err != nil {
		return &Config{}, err
	}
	return &conf, nil
}
