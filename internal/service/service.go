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

// Services struct holds all the service instances
type Services struct {
	Auth  AuthServiceInterface
	User  UserServiceInterface
	// Board BoardService
	// Column ColumnService
	// Card CardService
	// Comment CommentService
	// Label LabelService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth:  NewAuthService(repos.User, cfg),
		User:  NewUserService(repos.User),
		// Initialize other services here
	}
}