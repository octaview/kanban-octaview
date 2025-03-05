package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
	godotenv.Load()

	config := &Config{
		App: AppConfig{
			Env: getEnv("APP_ENV", "development"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "kanban"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	config.Database.DSN = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	jwtExpStr := getEnv("JWT_EXPIRATION", "24h")
	jwtExpiration, err := time.ParseDuration(jwtExpStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT expiration duration: %w", err)
	}

	config.JWT = JWTConfig{
		Secret:    getEnv("JWT_SECRET", "your_jwt_secret_key"),
		ExpiresIn: jwtExpiration,
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
