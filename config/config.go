package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RPCDialURL string `envconfig:"RPC_DIAL_URL" required:"true"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
