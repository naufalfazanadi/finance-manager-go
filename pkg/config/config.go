package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/naufalfazanadi/finance-manager-go/pkg/types"
)

type Config struct {
	Database types.DatabaseConfig
	Server   ServerConfig
	App      AppConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type AppConfig struct {
	Env      string
	LogLevel string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn string
}

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	return &Config{
		Database: types.DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "clean_api_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		App: AppConfig{
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "debug"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
