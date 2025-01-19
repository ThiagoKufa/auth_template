package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv carrega as variáveis de ambiente do arquivo .env
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("erro ao carregar arquivo .env: %v", err)
	}
	return nil
}

// GetEnv retorna o valor de uma variável de ambiente ou um valor padrão
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt retorna o valor de uma variável de ambiente como inteiro
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
