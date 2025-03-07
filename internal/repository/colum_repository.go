package repository

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/gorm"
)

type ColumnRepo struct {
	db *gorm.DB
}

func NewColumnRepo(db *gorm.DB) *ColumnRepo {
	return &ColumnRepo{db: db}
}

func (r *ColumnRepo) Create(ctx context.Context, column *models.Column) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var maxPosition struct {
			Max int
		}
		if err := tx.Model(&models.Column{}).
			Select("COALESCE(MAX(position), -1) as max").
			Where("board_id = ?", column.BoardID).
			Scan(&maxPosition).Error; err != nil {
			return models.NewDatabaseError("getting max column position", err)
		}

		column.Position = maxPosition.Max + 1

		if err := tx.Create(column).Error; err != nil {
			return models.NewDatabaseError("creating column", err)
		}

		return nil
	})
}

func (r *ColumnRepo) GetByID(ctx context.Context, id uint) (*models.Column, error) {
	var column models.Column
	result := r.db.WithContext(ctx).First(&column, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrColumnNotFound
		}
		return nil, models.NewDatabaseError("getting column by ID", result.Error)
	}
	return &column, nil
}

func (r *ColumnRepo) GetByBoardID(ctx context.Context, boardID uint) ([]models.Column, error) {
	var columns []models.Column
	result := r.db.WithContext(ctx).
		Where("board_id = ?", boardID).
		Order("position ASC").
		Find(&columns)
	if result.Error != nil {
		return nil, models.NewDatabaseError("getting columns by board ID", result.Error)
	}
	return columns, nil
}

func (r *ColumnRepo) Update(ctx context.Context, column *models.Column) error {
	result := r.db.WithContext(ctx).Save(column)
	if result.Error != nil {
		return models.NewDatabaseError("updating column", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrColumnNotFound
	}
	return nil
}

func (r *ColumnRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var column models.Column
		if err := tx.First(&column, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrColumnNotFound
			}
			return models.NewDatabaseError("finding column for deletion", err)
		}

		if err := tx.Delete(&column).Error; err != nil {
			return models.NewDatabaseError("deleting column", err)
		}

		if err := tx.Exec("UPDATE columns SET position = position - 1 WHERE board_id = ? AND position > ? AND deleted_at IS NULL",
			column.BoardID, column.Position).Error; err != nil {
			return models.NewDatabaseError("reordering columns after deletion", err)
		}

		return nil
	})
}

func (r *ColumnRepo) UpdatePositions(ctx context.Context, columns []models.Column) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, column := range columns {
			if err := tx.Model(&models.Column{}).
				Where("id = ?", column.ID).
				Update("position", i).Error; err != nil {
				return models.NewDatabaseError("updating column position", err)
			}
		}
		return nil
	})
}