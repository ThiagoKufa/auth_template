package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Log      LogConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port     string
	Timeout  time.Duration
	Compress bool
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type AuthConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

type SecurityConfig struct {
	CORS CORSConfig
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvStringSliceOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:     getEnvOrDefault("SERVER_PORT", "8081"),
			Timeout:  getEnvDurationOrDefault("SERVER_TIMEOUT", 30*time.Second),
			Compress: true,
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvIntOrDefault("DB_PORT", 5432),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			Name:     getEnvOrDefault("DB_NAME", "kufatech_dev"),
			SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     getEnvIntOrDefault("REDIS_PORT", 6379),
			Password: getEnvOrDefault("REDIS_PASSWORD", ""),
			DB:       getEnvIntOrDefault("REDIS_DB", 0),
			PoolSize: getEnvIntOrDefault("REDIS_POOL_SIZE", 10),
		},
		Auth: AuthConfig{
			AccessTokenSecret:  getEnvOrDefault("JWT_ACCESS_SECRET", "dev_access_secret"),
			RefreshTokenSecret: getEnvOrDefault("JWT_REFRESH_SECRET", "dev_refresh_secret"),
			AccessTokenTTL:     getEnvDurationOrDefault("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL:    getEnvDurationOrDefault("JWT_REFRESH_TTL", 720*time.Hour),
		},
		Log: LogConfig{
			Level:  getEnvOrDefault("LOG_LEVEL", "info"),
			Format: getEnvOrDefault("LOG_FORMAT", "json"),
		},
		Security: SecurityConfig{
			CORS: CORSConfig{
				AllowedOrigins:   getEnvStringSliceOrDefault("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:8080"}),
				AllowedMethods:   getEnvStringSliceOrDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
				AllowedHeaders:   getEnvStringSliceOrDefault("CORS_ALLOWED_HEADERS", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}),
				ExposedHeaders:   getEnvStringSliceOrDefault("CORS_EXPOSED_HEADERS", []string{"Link"}),
				AllowCredentials: true,
				MaxAge:           getEnvIntOrDefault("CORS_MAX_AGE", 86400),
			},
		},
	}

	return cfg, nil
}
