package models

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Position    int            `gorm:"not null" json:"position"`
	ColumnID    uint           `gorm:"not null" json:"column_id"`
	Column      Column         `gorm:"foreignKey:ColumnID" json:"column,omitempty"`
	AssignedTo  *uint          `json:"assigned_to,omitempty"`
	User        *User          `gorm:"foreignKey:AssignedTo" json:"user,omitempty"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
