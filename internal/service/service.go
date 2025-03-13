package service

import (
	"context"

	"github.com/octaview/kanban-octaview/internal/config"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/repository"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, user *models.User) (uint, error)
	Login(ctx context.Context, email, password string) (string, error)
	ParseToken(token string) (uint, error)
	RefreshToken(ctx context.Context, userID uint) (string, error)
}

type UserServiceInterface interface {
	Create(ctx context.Context, user *models.User) (uint, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error
}

type BoardServiceInterface interface {
	Create(ctx context.Context, board *models.Board) error
	GetByID(ctx context.Context, id uint) (*models.Board, error)
	GetByOwnerID(ctx context.Context, ownerID uint) ([]models.Board, error)
	Update(ctx context.Context, board *models.Board) error
	Delete(ctx context.Context, id uint) error
}

type ColumnServiceInterface interface {
	Create(ctx context.Context, column *models.Column) error
	GetByID(ctx context.Context, id uint) (*models.Column, error)
	GetByBoardID(ctx context.Context, boardID uint) ([]models.Column, error)
	Update(ctx context.Context, column *models.Column) error
	Delete(ctx context.Context, id uint) error
	UpdatePositions(ctx context.Context, columns []models.Column) error
}

// Services struct holds all the service instances
type Services struct {
	Auth   AuthServiceInterface
	User   UserServiceInterface
	Board  BoardServiceInterface
	Column ColumnServiceInterface
	// Card CardService
	// Comment CommentService
	// Label LabelService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth:   NewAuthService(repos.User, cfg),
		User:   NewUserService(repos.User),
		Board:  NewBoardService(repos.Board, repos.User),
		Column: NewColumnService(repos.Column, repos.Board),
		// Initialize other services here
	}
}