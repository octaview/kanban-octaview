package service

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/repository"
)

type ColumnService struct {
	columnRepo repository.ColumnRepository
	boardRepo  repository.BoardRepository
}

func NewColumnService(columnRepo repository.ColumnRepository, boardRepo repository.BoardRepository) *ColumnService {
	return &ColumnService{
		columnRepo: columnRepo,
		boardRepo:  boardRepo,
	}
}

func (s *ColumnService) Create(ctx context.Context, column *models.Column) error {
	_, err := s.boardRepo.GetByID(ctx, column.BoardID)
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			return models.ErrBoardNotFound
		}
		return err
	}

	return s.columnRepo.Create(ctx, column)
}

func (s *ColumnService) GetByID(ctx context.Context, id uint) (*models.Column, error) {
	return s.columnRepo.GetByID(ctx, id)
}

func (s *ColumnService) GetByBoardID(ctx context.Context, boardID uint) ([]models.Column, error) {
	_, err := s.boardRepo.GetByID(ctx, boardID)
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			return nil, models.ErrBoardNotFound
		}
		return nil, err
	}

	return s.columnRepo.GetByBoardID(ctx, boardID)
}

func (s *ColumnService) Update(ctx context.Context, column *models.Column) error {
	existingColumn, err := s.columnRepo.GetByID(ctx, column.ID)
	if err != nil {
		return err
	}

	if column.Position == 0 {
		column.Position = existingColumn.Position
	}

	if column.BoardID != 0 && column.BoardID != existingColumn.BoardID {
		return errors.New("cannot change board ID of a column")
	}

	if column.BoardID == 0 {
		column.BoardID = existingColumn.BoardID
	}

	return s.columnRepo.Update(ctx, column)
}

func (s *ColumnService) Delete(ctx context.Context, id uint) error {
	return s.columnRepo.Delete(ctx, id)
}

func (s *ColumnService) UpdatePositions(ctx context.Context, columns []models.Column) error {
	if len(columns) == 0 {
		return errors.New("no columns provided for position update")
	}

	boardID := columns[0].BoardID
	for _, col := range columns {
		if col.BoardID != boardID {
			return errors.New("all columns must belong to the same board")
		}
	}

	_, err := s.boardRepo.GetByID(ctx, boardID)
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			return models.ErrBoardNotFound
		}
		return err
	}

	return s.columnRepo.UpdatePositions(ctx, columns)
}
