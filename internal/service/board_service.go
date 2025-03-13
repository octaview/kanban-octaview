package service

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/repository"
)

type BoardService struct {
	repo     repository.BoardRepository
	userRepo repository.UserRepository
}

func NewBoardService(repo repository.BoardRepository, userRepo repository.UserRepository) *BoardService {
	return &BoardService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *BoardService) Create(ctx context.Context, board *models.Board) error {
	owner, err := s.userRepo.GetByID(ctx, board.OwnerID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.ErrUserNotFound
		}
		return err
	}
	
	board.OwnerID = owner.ID

	return s.repo.Create(ctx, board)
}

func (s *BoardService) GetByID(ctx context.Context, id uint) (*models.Board, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BoardService) GetByOwnerID(ctx context.Context, ownerID uint) ([]models.Board, error) {
	_, err := s.userRepo.GetByID(ctx, ownerID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return s.repo.GetByOwnerID(ctx, ownerID)
}

func (s *BoardService) Update(ctx context.Context, board *models.Board) error {
	existingBoard, err := s.repo.GetByID(ctx, board.ID)
	if err != nil {
		return err
	}

	if board.OwnerID == 0 {
		board.OwnerID = existingBoard.OwnerID
	} else if board.OwnerID != existingBoard.OwnerID {
		owner, err := s.userRepo.GetByID(ctx, board.OwnerID)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				return models.ErrUserNotFound
			}
			return err
		}
		board.OwnerID = owner.ID
	}

	return s.repo.Update(ctx, board)
}

func (s *BoardService) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}