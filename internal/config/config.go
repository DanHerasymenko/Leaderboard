package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env         string `env:"APP_ENV" envDefault:"local-local"`
	MongoURI    string `env:"MONGO_URI"`
	MongoDbName string `env:"MONGO_DB_NAME"`
	AuthSecret  string `env:"AUTH_SECRET"`
	ServerPort  string `env:"SERVER_ADDR" envDefault:":8082"`
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}
	return cfg, nil
}
