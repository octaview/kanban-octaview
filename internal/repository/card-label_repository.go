package repository

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/gorm"
)

type CardLabelRepo struct {
	db *gorm.DB
}

func NewCardLabelRepo(db *gorm.DB) *CardLabelRepo {
	return &CardLabelRepo{db: db}
}

func (r *CardLabelRepo) AddLabelToCard(ctx context.Context, cardID uint, labelID uint) error {
	cardLabel := models.CardLabel{
		CardID:  cardID,
		LabelID: labelID,
	}

	var existingCount int64
	if err := r.db.WithContext(ctx).Model(&models.CardLabel{}).
		Where("card_id = ? AND label_id = ?", cardID, labelID).
		Count(&existingCount).Error; err != nil {
		return models.NewDatabaseError("checking existing card-label association", err)
	}

	if existingCount > 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Create(&cardLabel)
	if result.Error != nil {
		return models.NewDatabaseError("adding label to card", result.Error)
	}
	return nil
}

func (r *CardLabelRepo) RemoveLabelFromCard(ctx context.Context, cardID uint, labelID uint) error {
	result := r.db.WithContext(ctx).
		Where("card_id = ? AND label_id = ?", cardID, labelID).
		Delete(&models.CardLabel{})
	
	if result.Error != nil {
		return models.NewDatabaseError("removing label from card", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return errors.New("card does not have this label")
	}
	
	return nil
}

func (r *CardLabelRepo) GetLabelsByCardID(ctx context.Context, cardID uint) ([]models.Label, error) {
	var labels []models.Label
	
	err := r.db.WithContext(ctx).
		Select("labels.*").
		Joins("JOIN card_labels ON labels.id = card_labels.label_id").
		Where("card_labels.card_id = ?", cardID).
		Find(&labels).Error
	
	if err != nil {
		return nil, models.NewDatabaseError("getting labels by card ID", err)
	}
	
	return labels, nil
}

func (r *CardLabelRepo) GetCardsByLabelID(ctx context.Context, labelID uint) ([]models.Card, error) {
	var cards []models.Card
	
	err := r.db.WithContext(ctx).
		Select("cards.*").
		Joins("JOIN card_labels ON cards.id = card_labels.card_id").
		Where("card_labels.label_id = ? AND cards.deleted_at IS NULL", labelID).
		Find(&cards).Error
	
	if err != nil {
		return nil, models.NewDatabaseError("getting cards by label ID", err)
	}
	
	return cards, nil
}