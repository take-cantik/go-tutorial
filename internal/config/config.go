package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv   string
	HTTPPort string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
}

func Load() Config {
	return Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "3306"),
		DBUser:   getEnv("DB_USER", "tutorial"),
		DBPass:   getEnv("DB_PASSWORD", "tutorial"),
		DBName:   getEnv("DB_NAME", "go_tutorial"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
