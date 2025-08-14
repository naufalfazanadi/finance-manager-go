package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string

	// Connection Pool Settings
	MaxOpenConns    int // Maximum open connections
	MaxIdleConns    int // Maximum idle connections
	ConnMaxLifetime int // Connection lifetime in minutes
	ConnMaxIdleTime int // Idle timeout in minutes

	// Retry Settings
	MaxRetries     int // Retry attempts
	RetryDelay     int // Delay between retries in seconds
	ConnectTimeout int // Initial connection timeout in seconds
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	App      AppConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Minio    MinioConfig
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

type CORSConfig struct {
	AllowOrigins string
	AllowMethods string
	AllowHeaders string
}

type MinioConfig struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	UseSSL        bool
	PrivateBucket string
	PublicBucket  string
	Directory     string
}

var globalConfig *Config

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	globalConfig = &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "clean_api_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),

			// Connection Pool Settings with defaults optimized for goroutines
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 30), // 30 minutes
			ConnMaxIdleTime: getEnvAsInt("DB_CONN_MAX_IDLE_TIME", 5), // 5 minutes

			// Retry Settings
			MaxRetries:     getEnvAsInt("DB_MAX_RETRIES", 3),      // 3 retry attempts
			RetryDelay:     getEnvAsInt("DB_RETRY_DELAY", 5),      // 5 seconds delay
			ConnectTimeout: getEnvAsInt("DB_CONNECT_TIMEOUT", 10), // 10 seconds timeout
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		App: AppConfig{
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "debug"),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		},
		CORS: CORSConfig{
			AllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
			AllowMethods: getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			AllowHeaders: getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		},
		Minio: MinioConfig{
			Endpoint:      getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:     getEnv("MINIO_SECRET_KEY", "minioadmin"),
			UseSSL:        getEnvAsBool("MINIO_USE_SSL", false),
			PrivateBucket: getEnv("MINIO_PRIVATE_BUCKET", "private"),
			PublicBucket:  getEnv("MINIO_PUBLIC_BUCKET", "public"),
			Directory:     getEnv("MINIO_DIRECTORY", ""),
		},
	}

	return globalConfig
}

func GetConfig() *Config {
	if globalConfig == nil {
		return LoadConfig()
	}
	return globalConfig
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using fallback: %d", key, value, fallback)
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
		log.Printf("Warning: Invalid boolean value for %s: %s, using fallback: %t", key, value, fallback)
	}
	return fallback
}
