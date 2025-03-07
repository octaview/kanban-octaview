package repository

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/gorm"
)

type BoardRepo struct {
	db *gorm.DB
}

func NewBoardRepo(db *gorm.DB) *BoardRepo {
	return &BoardRepo{db: db}
}

func (r *BoardRepo) Create(ctx context.Context, board *models.Board) error {
	result := r.db.WithContext(ctx).Create(board)
	if result.Error != nil {
		return models.NewDatabaseError("creating board", result.Error)
	}
	return nil
}

func (r *BoardRepo) GetByID(ctx context.Context, id uint) (*models.Board, error) {
	var board models.Board
	result := r.db.WithContext(ctx).First(&board, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrBoardNotFound
		}
		return nil, models.NewDatabaseError("getting board by ID", result.Error)
	}
	return &board, nil
}

func (r *BoardRepo) GetByOwnerID(ctx context.Context, ownerID uint) ([]models.Board, error) {
	var boards []models.Board
	result := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&boards)
	if result.Error != nil {
		return nil, models.NewDatabaseError("getting boards by owner ID", result.Error)
	}
	return boards, nil
}

func (r *BoardRepo) Update(ctx context.Context, board *models.Board) error {
	result := r.db.WithContext(ctx).Save(board)
	if result.Error != nil {
		return models.NewDatabaseError("updating board", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrBoardNotFound
	}
	return nil
}

func (r *BoardRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var board models.Board
		if err := tx.First(&board, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrBoardNotFound
			}
			return models.NewDatabaseError("finding board for deletion", err)
		}

		if err := tx.Delete(&board).Error; err != nil {
			return models.NewDatabaseError("deleting board", err)
		}

		return nil
	})
}