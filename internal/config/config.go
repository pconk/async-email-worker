// config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	RedisAddress string
	ApiPort      string
	QueueName    string
	WorkerNumber int
	MaxRetry     int
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	// Cek apakah string bisa diubah ke int
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}

func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()

	cfg := &AppConfig{
		RedisAddress: os.Getenv("REDIS_ADDRESS"),
		ApiPort:      os.Getenv("APP_PORT"),
		QueueName:    os.Getenv("QUEUE_NAME"),
		WorkerNumber: getEnvInt("WORKER_CONCURRENCY", 5),
		MaxRetry:     getEnvInt("MAX_RETRY", 3),
	}

	log.Println("cfg", cfg)
	if cfg.RedisAddress == "" {
		return nil, fmt.Errorf("REDIS_ADDRESS is required")
	}
	if cfg.ApiPort == "" {
		cfg.ApiPort = "8081"
	}
	return cfg, nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
