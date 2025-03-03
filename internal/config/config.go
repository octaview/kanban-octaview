package config

import (
	"time"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Env string
}

type HTTPConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}