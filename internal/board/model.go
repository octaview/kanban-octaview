package board

import "time"

type Board struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type Column struct {
	ID      int64  `json:"id"`
	BoardID int64  `json:"board_id"`
	Title   string `json:"title" validate:"required"`
}

type Task struct {
	ID        int64     `json:"id"`
	ColumnID  int64     `json:"column_id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	DueDate   time.Time `json:"due_date"`
}
