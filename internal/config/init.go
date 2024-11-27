package config

import (
	"encoding/json"
	"github.com/caarlos0/env"
	"log"
	"os"
)

func NewConfig() (*Config, error) {
	cfg := Config{}
	data, err := os.ReadFile("./configs/config.json")
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err = env.Parse(&cfg); err != nil {
		return nil, err
	}
	log.Println(cfg)
	return &cfg, nil
}
