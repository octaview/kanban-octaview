package models

type Label struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Color   string `gorm:"not null" json:"color"`
	BoardID uint   `gorm:"not null" json:"board_id"`
	Board   Board  `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}
