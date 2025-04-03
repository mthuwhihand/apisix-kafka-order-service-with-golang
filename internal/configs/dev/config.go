package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	instance *Config
	once     sync.Once
	logger   = log.New(log.Writer(), "[configs/dev/config.go] ", log.LstdFlags|log.Lshortfile)
)

type Config struct {
	DB_URL    string
	SecretKey string
}

func Instance() (*Config, error) {
	var err error
	rootDir, err := os.Getwd()
	if err != nil {
		logger.Println("Error getting working directory:", err)
	} else {
		logger.Println("Current working directory: ", rootDir)
	}

	once.Do(func() {
		err = godotenv.Load("internal/configs/dev/.env")
		if err != nil {
			logger.Println("Can not load env file!")
		} else {
			logger.Println("Env file loaded successfully")
		}

		instance = &Config{
			DB_URL:    os.Getenv("DB_URL"),
			SecretKey: os.Getenv("SECRET_KEY"),
		}
	})

	if instance.DB_URL == "" || instance.SecretKey == "" {
		logger.Println("Some env variable may be missing")
	} else {
		logger.Println("Loaded environment variables: DB_URL = ", instance.DB_URL)
	}

	return instance, err
}

func (cfg *Config) GetSecretKey() string {
	return cfg.SecretKey
}
