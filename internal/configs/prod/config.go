package config

import (
	"log"
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	Mongo     MongoConfig
	SecretKey string
}

type MongoConfig struct {
	URI      string
	Database string
}

func Instance() (*Config, error) {
	var err error
	once.Do(func() {

		instance = &Config{}
		log.Println("unable to decode into struct")
	})

	return instance, err
}

// GetSecretKey trả về SecretKey
func (cfg *Config) GetSecretKey() string {
	return cfg.SecretKey
}
