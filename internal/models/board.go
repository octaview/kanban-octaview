package models

import (
	"time"

	"gorm.io/gorm"
)

type Board struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	OwnerID     uint           `gorm:"not null" json:"owner_id"`
	Owner       User           `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
