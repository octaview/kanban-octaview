package service

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/repository"
)

type LabelService struct {
	labelRepo repository.LabelRepository
	boardRepo repository.BoardRepository
}

func NewLabelService(labelRepo repository.LabelRepository, boardRepo repository.BoardRepository) *LabelService {
	return &LabelService{
		labelRepo: labelRepo,
		boardRepo: boardRepo,
	}
}

func (s *LabelService) Create(ctx context.Context, label *models.Label) error {
	_, err := s.boardRepo.GetByID(ctx, label.BoardID)
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			return models.ErrBoardNotFound
		}
		return err
	}

	if label.Color == "" {
		return models.NewValidationError("color", "color is required")
	}

	if label.Name == "" {
		return models.NewValidationError("name", "name is required")
	}

	return s.labelRepo.Create(ctx, label)
}

func (s *LabelService) GetByID(ctx context.Context, id uint) (*models.Label, error) {
	return s.labelRepo.GetByID(ctx, id)
}

func (s *LabelService) GetByBoardID(ctx context.Context, boardID uint) ([]models.Label, error) {
	_, err := s.boardRepo.GetByID(ctx, boardID)
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			return nil, models.ErrBoardNotFound
		}
		return nil, err
	}

	return s.labelRepo.GetByBoardID(ctx, boardID)
}

func (s *LabelService) Update(ctx context.Context, label *models.Label) error {
	existingLabel, err := s.labelRepo.GetByID(ctx, label.ID)
	if err != nil {
		return err
	}

	if label.BoardID != 0 && label.BoardID != existingLabel.BoardID {
		return errors.New("cannot change board ID of a label")
	}

	if label.BoardID == 0 {
		label.BoardID = existingLabel.BoardID
	}

	if label.Name == "" {
		return models.NewValidationError("name", "name is required")
	}

	if label.Color == "" {
		return models.NewValidationError("color", "color is required")
	}

	return s.labelRepo.Update(ctx, label)
}

func (s *LabelService) Delete(ctx context.Context, id uint) error {
	_, err := s.labelRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.labelRepo.Delete(ctx, id)
}