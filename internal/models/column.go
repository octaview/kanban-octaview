package models

import (
	"time"

	"gorm.io/gorm"
)

type Column struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Position  int            `gorm:"not null" json:"position"`
	BoardID   uint           `gorm:"not null" json:"board_id"`
	Board     Board          `gorm:"foreignKey:BoardID" json:"board,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
