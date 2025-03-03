package database

import (
	"fmt"
	"log"

	"github.com/octaview/kanban-octaview/internal/config"
	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.App.Env == "development" {
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if cfg.App.Env == "development" {
		err = db.AutoMigrate(
			&models.User{},
			&models.Board{},
			&models.Column{},
			&models.Card{},
			&models.Label{},
			&models.Comment{},
		)
		if err != nil {
			log.Printf("Warning: Auto migration failed: %v", err)
		}
	}

	return db, nil
}
