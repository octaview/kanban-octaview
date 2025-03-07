package repository

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/gorm"
)

type CardRepo struct {
	db *gorm.DB
}

func NewCardRepo(db *gorm.DB) *CardRepo {
	return &CardRepo{db: db}
}

func (r *CardRepo) Create(ctx context.Context, card *models.Card) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var maxPosition struct {
			Max int
		}
		if err := tx.Model(&models.Card{}).
			Select("COALESCE(MAX(position), -1) as max").
			Where("column_id = ?", card.ColumnID).
			Scan(&maxPosition).Error; err != nil {
			return models.NewDatabaseError("getting max card position", err)
		}

		card.Position = maxPosition.Max + 1

		if err := tx.Create(card).Error; err != nil {
			return models.NewDatabaseError("creating card", err)
		}

		return nil
	})
}

func (r *CardRepo) GetByID(ctx context.Context, id uint) (*models.Card, error) {
	var card models.Card
	result := r.db.WithContext(ctx).First(&card, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrCardNotFound
		}
		return nil, models.NewDatabaseError("getting card by ID", result.Error)
	}
	return &card, nil
}

func (r *CardRepo) GetByColumnID(ctx context.Context, columnID uint) ([]models.Card, error) {
	var cards []models.Card
	result := r.db.WithContext(ctx).
		Where("column_id = ?", columnID).
		Order("position ASC").
		Find(&cards)
	if result.Error != nil {
		return nil, models.NewDatabaseError("getting cards by column ID", result.Error)
	}
	return cards, nil
}

func (r *CardRepo) Update(ctx context.Context, card *models.Card) error {
	result := r.db.WithContext(ctx).Save(card)
	if result.Error != nil {
		return models.NewDatabaseError("updating card", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrCardNotFound
	}
	return nil
}

func (r *CardRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var card models.Card
		if err := tx.First(&card, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrCardNotFound
			}
			return models.NewDatabaseError("finding card for deletion", err)
		}

		if err := tx.Delete(&card).Error; err != nil {
			return models.NewDatabaseError("deleting card", err)
		}

		if err := tx.Exec("UPDATE cards SET position = position - 1 WHERE column_id = ? AND position > ? AND deleted_at IS NULL",
			card.ColumnID, card.Position).Error; err != nil {
			return models.NewDatabaseError("reordering cards after deletion", err)
		}

		return nil
	})
}

func (r *CardRepo) UpdatePositions(ctx context.Context, cards []models.Card) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, card := range cards {
			if err := tx.Model(&models.Card{}).
				Where("id = ?", card.ID).
				Update("position", i).Error; err != nil {
				return models.NewDatabaseError("updating card position", err)
			}
		}
		return nil
	})
}

func (r *CardRepo) MoveToColumn(ctx context.Context, cardID, columnID uint, position int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var card models.Card
		if err := tx.First(&card, cardID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ErrCardNotFound
			}
			return models.NewDatabaseError("finding card for moving", err)
		}

		oldColumnID := card.ColumnID
		oldPosition := card.Position

		if oldColumnID == columnID {
			if position < oldPosition {
				if err := tx.Exec("UPDATE cards SET position = position + 1 WHERE column_id = ? AND position >= ? AND position < ? AND deleted_at IS NULL",
					columnID, position, oldPosition).Error; err != nil {
					return models.NewDatabaseError("shifting cards for move within column (up)", err)
				}
			}
			if position > oldPosition {
				if err := tx.Exec("UPDATE cards SET position = position - 1 WHERE column_id = ? AND position > ? AND position <= ? AND deleted_at IS NULL",
					columnID, oldPosition, position).Error; err != nil {
					return models.NewDatabaseError("shifting cards for move within column (down)", err)
				}
			}
		} else {

			if err := tx.Exec("UPDATE cards SET position = position + 1 WHERE column_id = ? AND position >= ? AND deleted_at IS NULL",
				columnID, position).Error; err != nil {
				return models.NewDatabaseError("shifting cards in new column", err)
			}

			if err := tx.Exec("UPDATE cards SET position = position - 1 WHERE column_id = ? AND position > ? AND deleted_at IS NULL",
				oldColumnID, oldPosition).Error; err != nil {
				return models.NewDatabaseError("shifting cards in old column", err)
			}
		}

		if err := tx.Model(&card).Updates(map[string]interface{}{
			"column_id": columnID,
			"position":  position,
		}).Error; err != nil {
			return models.NewDatabaseError("updating card column and position", err)
		}

		return nil
	})
}