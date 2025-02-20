package repository

import (
	"context"
	"github.com/yourusername/kanban-backend/internal/board"
)

type BoardRepository interface {
	Create(ctx context.Context, board *board.Board) (int64, error)
	GetByID(ctx context.Context, id int64) (*board.Board, error)
	// Дополнительные методы: Update, Delete, List и пр.
}
