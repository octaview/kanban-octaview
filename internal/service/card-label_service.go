package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/repository"
)

type CardLabelRepository interface {
	AddLabelToCard(ctx context.Context, cardID uint, labelID uint) error
	RemoveLabelFromCard(ctx context.Context, cardID uint, labelID uint) error
	GetLabelsByCardID(ctx context.Context, cardID uint) ([]models.Label, error)
	GetCardsByLabelID(ctx context.Context, labelID uint) ([]models.Card, error)
}

type CardLabelService struct {
	cardLabelRepo CardLabelRepository
	cardRepo      repository.CardRepository
	labelRepo     repository.LabelRepository
	boardRepo     repository.BoardRepository
	columnRepo    repository.ColumnRepository
}

func NewCardLabelService(
	cardLabelRepo CardLabelRepository,
	cardRepo repository.CardRepository,
	labelRepo repository.LabelRepository,
	boardRepo repository.BoardRepository,
	columnRepo repository.ColumnRepository,
) *CardLabelService {
	return &CardLabelService{
		cardLabelRepo: cardLabelRepo,
		cardRepo:      cardRepo,
		labelRepo:     labelRepo,
		boardRepo:     boardRepo,
		columnRepo:    columnRepo,
	}
}

func (s *CardLabelService) AddLabelToCard(ctx context.Context, cardID uint, labelID uint) error {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, models.ErrCardNotFound) {
			return models.ErrCardNotFound
		}
		return err
	}

	label, err := s.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		if errors.Is(err, models.ErrLabelNotFound) {
			return models.ErrLabelNotFound
		}
		return err
	}

	column, err := s.getColumnForCard(ctx, card)
	if err != nil {
		return err
	}

	if label.BoardID != column.BoardID {
		return fmt.Errorf("label %d belongs to board %d but card %d belongs to board %d", 
			labelID, label.BoardID, cardID, column.BoardID)
	}

	return s.cardLabelRepo.AddLabelToCard(ctx, cardID, labelID)
}

func (s *CardLabelService) RemoveLabelFromCard(ctx context.Context, cardID uint, labelID uint) error {
	_, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, models.ErrCardNotFound) {
			return models.ErrCardNotFound
		}
		return err
	}

	_, err = s.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		if errors.Is(err, models.ErrLabelNotFound) {
			return models.ErrLabelNotFound
		}
		return err
	}

	return s.cardLabelRepo.RemoveLabelFromCard(ctx, cardID, labelID)
}

func (s *CardLabelService) GetLabelsByCardID(ctx context.Context, cardID uint) ([]models.Label, error) {
	_, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, models.ErrCardNotFound) {
			return nil, models.ErrCardNotFound
		}
		return nil, err
	}

	return s.cardLabelRepo.GetLabelsByCardID(ctx, cardID)
}

func (s *CardLabelService) GetCardsByLabelID(ctx context.Context, labelID uint) ([]models.Card, error) {
	_, err := s.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		if errors.Is(err, models.ErrLabelNotFound) {
			return nil, models.ErrLabelNotFound
		}
		return nil, err
	}

	return s.cardLabelRepo.GetCardsByLabelID(ctx, labelID)
}

func (s *CardLabelService) getColumnForCard(ctx context.Context, card *models.Card) (*models.Column, error) {
	if card.Column.ID != 0 {
		return &card.Column, nil
	}
	
	column, err := s.getColumnByID(ctx, card.ColumnID)
	if err != nil {
		return nil, err
	}
	
	return column, nil
}

func (s *CardLabelService) getColumnByID(ctx context.Context, columnID uint) (*models.Column, error) {
	column, err := s.columnRepo.GetByID(ctx, columnID)
	if err != nil {
		if errors.Is(err, models.ErrColumnNotFound) {
			return nil, models.ErrColumnNotFound
		}
		return nil, fmt.Errorf("failed to get column: %w", err)
	}
	
	return column, nil
}

func (s *CardLabelService) BatchAddLabelsToCard(ctx context.Context, cardID uint, labelIDs []uint) error {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, models.ErrCardNotFound) {
			return models.ErrCardNotFound
		}
		return err
	}
	
	column, err := s.getColumnForCard(ctx, card)
	if err != nil {
		return err
	}
	
	for _, labelID := range labelIDs {
		label, err := s.labelRepo.GetByID(ctx, labelID)
		if err != nil {
			if errors.Is(err, models.ErrLabelNotFound) {
				return fmt.Errorf("label %d not found", labelID)
			}
			return err
		}
		
		if label.BoardID != column.BoardID {
			return fmt.Errorf("label %d belongs to board %d but card %d belongs to board %d", 
				labelID, label.BoardID, cardID, column.BoardID)
		}
	}
	
	for _, labelID := range labelIDs {
		if err := s.cardLabelRepo.AddLabelToCard(ctx, cardID, labelID); err != nil {
			return err
		}
	}
	
	return nil
}

func (s *CardLabelService) BatchRemoveLabelsFromCard(ctx context.Context, cardID uint, labelIDs []uint) error {
	_, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, models.ErrCardNotFound) {
			return models.ErrCardNotFound
		}
		return err
	}
	
	for _, labelID := range labelIDs {
		_, err := s.labelRepo.GetByID(ctx, labelID)
		if err != nil {
			if errors.Is(err, models.ErrLabelNotFound) {
				return fmt.Errorf("label %d not found", labelID)
			}
			return err
		}
	}
	
	for _, labelID := range labelIDs {
		if err := s.cardLabelRepo.RemoveLabelFromCard(ctx, cardID, labelID); err != nil {
			return err
		}
	}
	
	return nil
}

func (s *CardLabelService) GetCardCountByLabelID(ctx context.Context, labelID uint) (int, error) {
	cards, err := s.GetCardsByLabelID(ctx, labelID)
	if err != nil {
		return 0, err
	}
	
	return len(cards), nil
}

func (s *CardLabelService) GetLabelCountByCardID(ctx context.Context, cardID uint) (int, error) {
	labels, err := s.GetLabelsByCardID(ctx, cardID)
	if err != nil {
		return 0, err
	}
	
	return len(labels), nil
}

func (s *CardLabelService) RemoveAllLabelsFromCard(ctx context.Context, cardID uint) error {
	_, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, models.ErrCardNotFound) {
			return models.ErrCardNotFound
		}
		return err
	}
	
	labels, err := s.GetLabelsByCardID(ctx, cardID)
	if err != nil {
		return err
	}
	
	for _, label := range labels {
		if err := s.cardLabelRepo.RemoveLabelFromCard(ctx, cardID, label.ID); err != nil {
			return err
		}
	}
	
	return nil
}