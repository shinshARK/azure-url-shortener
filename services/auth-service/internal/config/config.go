package config

import (
	"os"
)

type Config struct {
	Port         string
	DBHost       string
	DBName       string
	DBUser       string
	DBPassword   string
	JWTSecret    string
}

func LoadConfig() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBName:     getEnv("DB_NAME", "UrlShortenerDb"),
		DBUser:     getEnv("DB_USER", "sa"),
		DBPassword: getEnv("DB_PASSWORD", "yourStrong(!)Password"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
