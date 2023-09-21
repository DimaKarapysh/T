package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	HTTPAddr string
}

func Read() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	httpAddr, exists := os.LookupEnv("HTTP_ADDR")
	if exists {
		config.HTTPAddr = httpAddr
	}
	return config, nil
}
