package models

import "time"

type CardLabel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CardID    uint      `gorm:"not null" json:"card_id"`
	Card      Card      `gorm:"foreignKey:CardID" json:"card,omitempty"`
	LabelID   uint      `gorm:"not null" json:"label_id"`
	Label     Label     `gorm:"foreignKey:LabelID" json:"label,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}