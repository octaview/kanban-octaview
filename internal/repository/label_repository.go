package repository

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/gorm"
)

type LabelRepo struct {
	db *gorm.DB
}

func NewLabelRepo(db *gorm.DB) *LabelRepo {
	return &LabelRepo{db: db}
}

func (r *LabelRepo) Create(ctx context.Context, label *models.Label) error {
	result := r.db.WithContext(ctx).Create(label)
	if result.Error != nil {
		return models.NewDatabaseError("creating label", result.Error)
	}
	return nil
}

func (r *LabelRepo) GetByID(ctx context.Context, id uint) (*models.Label, error) {
	var label models.Label
	result := r.db.WithContext(ctx).First(&label, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrLabelNotFound
		}
		return nil, models.NewDatabaseError("getting label by ID", result.Error)
	}
	return &label, nil
}

func (r *LabelRepo) GetByBoardID(ctx context.Context, boardID uint) ([]models.Label, error) {
	var labels []models.Label
	result := r.db.WithContext(ctx).
		Where("board_id = ?", boardID).
		Find(&labels)
	if result.Error != nil {
		return nil, models.NewDatabaseError("getting labels by board ID", result.Error)
	}
	return labels, nil
}

func (r *LabelRepo) Update(ctx context.Context, label *models.Label) error {
	result := r.db.WithContext(ctx).Save(label)
	if result.Error != nil {
		return models.NewDatabaseError("updating label", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrLabelNotFound
	}
	return nil
}

func (r *LabelRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Label{}, id)
	if result.Error != nil {
		return models.NewDatabaseError("deleting label", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrLabelNotFound
	}
	return nil
}