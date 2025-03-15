package service

import (
	"context"
	"errors"
	"time"

	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/repository"
)

type CardService struct {
	cardRepo   repository.CardRepository
	columnRepo repository.ColumnRepository
	userRepo   repository.UserRepository
}

func NewCardService(cardRepo repository.CardRepository, columnRepo repository.ColumnRepository, userRepo repository.UserRepository) *CardService {
	return &CardService{
		cardRepo:   cardRepo,
		columnRepo: columnRepo,
		userRepo:   userRepo,
	}
}

func (s *CardService) Create(ctx context.Context, card *models.Card) error {
	_, err := s.columnRepo.GetByID(ctx, card.ColumnID)
	if err != nil {
		if errors.Is(err, models.ErrColumnNotFound) {
			return models.ErrColumnNotFound
		}
		return err
	}

	if card.AssignedTo != nil {
		_, err := s.userRepo.GetByID(ctx, *card.AssignedTo)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				return models.NewValidationError("assigned_to", "user not found")
			}
			return err
		}
	}

	if card.DueDate != nil && card.DueDate.Before(time.Now()) {
		return models.NewValidationError("due_date", "due date cannot be in the past")
	}

	return s.cardRepo.Create(ctx, card)
}

func (s *CardService) GetByID(ctx context.Context, id uint) (*models.Card, error) {
	return s.cardRepo.GetByID(ctx, id)
}

func (s *CardService) GetByColumnID(ctx context.Context, columnID uint) ([]models.Card, error) {
	_, err := s.columnRepo.GetByID(ctx, columnID)
	if err != nil {
		if errors.Is(err, models.ErrColumnNotFound) {
			return nil, models.ErrColumnNotFound
		}
		return nil, err
	}

	return s.cardRepo.GetByColumnID(ctx, columnID)
}

func (s *CardService) Update(ctx context.Context, card *models.Card) error {
	existingCard, err := s.cardRepo.GetByID(ctx, card.ID)
	if err != nil {
		return err
	}

	if card.ColumnID != 0 && card.ColumnID != existingCard.ColumnID {
		_, err := s.columnRepo.GetByID(ctx, card.ColumnID)
		if err != nil {
			if errors.Is(err, models.ErrColumnNotFound) {
				return models.ErrColumnNotFound
			}
			return err
		}
	} else {
		card.ColumnID = existingCard.ColumnID
	}

	if card.AssignedTo != nil && (existingCard.AssignedTo == nil || *card.AssignedTo != *existingCard.AssignedTo) {
		_, err := s.userRepo.GetByID(ctx, *card.AssignedTo)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				return models.NewValidationError("assigned_to", "user not found")
			}
			return err
		}
	}

	if card.DueDate != nil && (existingCard.DueDate == nil || !card.DueDate.Equal(*existingCard.DueDate)) {
		if card.DueDate.Before(time.Now()) {
			return models.NewValidationError("due_date", "due date cannot be in the past")
		}
	}

	if card.Position == 0 {
		card.Position = existingCard.Position
	}

	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) Delete(ctx context.Context, id uint) error {
	_, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.cardRepo.Delete(ctx, id)
}

func (s *CardService) UpdatePositions(ctx context.Context, cards []models.Card) error {
	if len(cards) == 0 {
		return errors.New("no cards provided for position update")
	}

	columnID := cards[0].ColumnID
	for _, card := range cards {
		if card.ColumnID != columnID {
			return errors.New("all cards must belong to the same column")
		}
	}

	_, err := s.columnRepo.GetByID(ctx, columnID)
	if err != nil {
		if errors.Is(err, models.ErrColumnNotFound) {
			return models.ErrColumnNotFound
		}
		return err
	}

	return s.cardRepo.UpdatePositions(ctx, cards)
}

func (s *CardService) MoveCard(ctx context.Context, cardID, columnID uint, position int) error {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return err
	}

	_, err = s.columnRepo.GetByID(ctx, columnID)
	if err != nil {
		if errors.Is(err, models.ErrColumnNotFound) {
			return models.ErrColumnNotFound
		}
		return err
	}

	if card.ColumnID == columnID && card.Position == position {
		return nil
	}

	return s.cardRepo.MoveToColumn(ctx, cardID, columnID, position)
}

func (s *CardService) AssignCard(ctx context.Context, cardID, userID uint) error {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return err
	}

	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.ErrUserNotFound
		}
		return err
	}

	card.AssignedTo = &userID
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) UnassignCard(ctx context.Context, cardID uint) error {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return err
	}

	card.AssignedTo = nil
	return s.cardRepo.Update(ctx, card)
}

func (s *CardService) UpdateDueDate(ctx context.Context, cardID uint, dueDate *time.Time) error {
	card, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return err
	}

	if dueDate != nil && dueDate.Before(time.Now()) {
		return models.NewValidationError("due_date", "due date cannot be in the past")
	}

	card.DueDate = dueDate
	return s.cardRepo.Update(ctx, card)
}