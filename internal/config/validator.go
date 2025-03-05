package config

import (
	"errors"
	"net"
	"slices"
	"strconv"
	"strings"
)

func (c *Config) Validate() error {
	if err := validateAppConfig(c.App); err != nil {
		return err
	}

	if err := validateHTTPConfig(c.HTTP); err != nil {
		return err
	}

	if err := validateDatabaseConfig(c.Database); err != nil {
		return err
	}

	if err := validateJWTConfig(c.JWT); err != nil {
		return err
	}

	return nil
}

func validateAppConfig(app AppConfig) error {
	validEnvs := []string{"development", "production", "staging", "test"}
	if !slices.Contains(validEnvs, app.Env) {
		return errors.New("invalid APP_ENV: must be one of development, production, staging, or test")
	}
	return nil
}

func validateHTTPConfig(http HTTPConfig) error {
	port, err := strconv.Atoi(http.Port)
	if err != nil {
		return errors.New("invalid HTTP_PORT: must be a valid integer")
	}
	if port < 1 || port > 65535 {
		return errors.New("invalid HTTP_PORT: must be between 1 and 65535")
	}
	return nil
}

func validateDatabaseConfig(db DatabaseConfig) error {
	if strings.TrimSpace(db.Host) == "" {
		return errors.New("DB_HOST cannot be empty")
	}

	port, err := strconv.Atoi(db.Port)
	if err != nil {
		return errors.New("invalid DB_PORT: must be a valid integer")
	}
	if port < 1 || port > 65535 {
		return errors.New("invalid DB_PORT: must be between 1 and 65535")
	}

	if strings.TrimSpace(db.User) == "" {
		return errors.New("DB_USER cannot be empty")
	}

	if strings.TrimSpace(db.DBName) == "" {
		return errors.New("DB_NAME cannot be empty")
	}

	validSSLModes := []string{"disable", "allow", "prefer", "require", "verify-full", "verify-ca"}
	if !slices.Contains(validSSLModes, db.SSLMode) {
		return errors.New("invalid DB_SSLMODE: must be one of disable, require, verify-full, or verify-ca")
	}

	if host := db.Host; host != "localhost" {
		_, err := net.LookupIP(host)
		if err != nil {
			return errors.New("invalid DB_HOST: unable to resolve hostname")
		}
	}

	return nil
}

func validateJWTConfig(jwt JWTConfig) error {
	if strings.TrimSpace(jwt.Secret) == "" {
		return errors.New("JWT_SECRET cannot be empty")
	}

	if jwt.ExpiresIn <= 0 {
		return errors.New("JWT_EXPIRATION must be a positive duration")
	}

	return nil
}