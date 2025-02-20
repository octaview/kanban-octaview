package repository

import (
	"context"
	"database/sql"
	"kanban-octaview/internal/board"
)

type boardRepo struct {
	db *sql.DB
}

func NewBoardRepository(db *sql.DB) BoardRepository {
	return &boardRepo{db: db}
}

func (r *boardRepo) Create(ctx context.Context, b *board.Board) (int64, error) {
	var id int64
	query := `INSERT INTO boards (title) VALUES ($1) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, b.Title).Scan(&id)
	return id, err
}

func (r *boardRepo) GetByID(ctx context.Context, id int64) (*board.Board, error) {
	b := &board.Board{}
	query := `SELECT id, title, created_at FROM boards WHERE id=$1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.Title, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return b, err
}
