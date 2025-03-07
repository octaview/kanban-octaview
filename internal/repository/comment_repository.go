package repository

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *models.Comment) error {
	result := r.db.WithContext(ctx).Create(comment)
	if result.Error != nil {
		return models.NewDatabaseError("creating comment", result.Error)
	}
	return nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id uint) (*models.Comment, error) {
	var comment models.Comment
	result := r.db.WithContext(ctx).First(&comment, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrCommentNotFound
		}
		return nil, models.NewDatabaseError("getting comment by ID", result.Error)
	}
	return &comment, nil
}

func (r *CommentRepo) GetByCardID(ctx context.Context, cardID uint) ([]models.Comment, error) {
	var comments []models.Comment
	result := r.db.WithContext(ctx).
		Where("card_id = ?", cardID).
		Order("created_at DESC").
		Find(&comments)
	if result.Error != nil {
		return nil, models.NewDatabaseError("getting comments by card ID", result.Error)
	}
	return comments, nil
}

func (r *CommentRepo) Update(ctx context.Context, comment *models.Comment) error {
	result := r.db.WithContext(ctx).Save(comment)
	if result.Error != nil {
		return models.NewDatabaseError("updating comment", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrCommentNotFound
	}
	return nil
}

func (r *CommentRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Comment{}, id)
	if result.Error != nil {
		return models.NewDatabaseError("deleting comment", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrCommentNotFound
	}
	return nil
}