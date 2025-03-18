package service

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/repository"
)

type CommentService struct {
	commentRepo repository.CommentRepository
	cardRepo    repository.CardRepository
	userRepo    repository.UserRepository
}

func NewCommentService(commentRepo repository.CommentRepository, cardRepo repository.CardRepository, userRepo repository.UserRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		cardRepo:    cardRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) Create(ctx context.Context, comment *models.Comment) error {
	card, err := s.cardRepo.GetByID(ctx, comment.CardID)
	if err != nil {
		if errors.Is(err, models.ErrCardNotFound) {
			return models.ErrCardNotFound
		}
		return err
	}

	user, err := s.userRepo.GetByID(ctx, comment.UserID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return models.ErrUserNotFound
		}
		return err
	}

	comment.Card = *card
	comment.User = *user
	
	return s.commentRepo.Create(ctx, comment)
}

func (s *CommentService) GetByID(ctx context.Context, id uint) (*models.Comment, error) {
	return s.commentRepo.GetByID(ctx, id)
}

func (s *CommentService) GetByCardID(ctx context.Context, cardID uint) ([]models.Comment, error) {
	_, err := s.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, models.ErrCardNotFound) {
			return nil, models.ErrCardNotFound
		}
		return nil, err
	}

	return s.commentRepo.GetByCardID(ctx, cardID)
}

func (s *CommentService) Update(ctx context.Context, comment *models.Comment) error {
	existingComment, err := s.commentRepo.GetByID(ctx, comment.ID)
	if err != nil {
		return err
	}

	existingComment.Content = comment.Content
	
	return s.commentRepo.Update(ctx, existingComment)
}

func (s *CommentService) Delete(ctx context.Context, id uint) error {
	_, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	return s.commentRepo.Delete(ctx, id)
}