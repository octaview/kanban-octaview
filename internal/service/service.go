package service

import (
	"context"
	"time"

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

type CardServiceInterface interface {
	Create(ctx context.Context, card *models.Card) error
	GetByID(ctx context.Context, id uint) (*models.Card, error)
	GetByColumnID(ctx context.Context, columnID uint) ([]models.Card, error)
	Update(ctx context.Context, card *models.Card) error
	Delete(ctx context.Context, id uint) error
	UpdatePositions(ctx context.Context, cards []models.Card) error
	MoveCard(ctx context.Context, cardID, columnID uint, position int) error
	AssignCard(ctx context.Context, cardID, userID uint) error
	UnassignCard(ctx context.Context, cardID uint) error
	UpdateDueDate(ctx context.Context, cardID uint, dueDate *time.Time) error
}

type LabelServiceInterface interface {
	Create(ctx context.Context, label *models.Label) error
	GetByID(ctx context.Context, id uint) (*models.Label, error)
	GetByBoardID(ctx context.Context, boardID uint) ([]models.Label, error)
	Update(ctx context.Context, label *models.Label) error
	Delete(ctx context.Context, id uint) error
}

// Services struct holds all the service instances
type Services struct {
	Auth   AuthServiceInterface
	User   UserServiceInterface
	Board  BoardServiceInterface
	Column ColumnServiceInterface
	Card   CardServiceInterface
	Label  LabelServiceInterface
	// Comment CommentService
}

func NewServices(repos *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth:   NewAuthService(repos.User, cfg),
		User:   NewUserService(repos.User),
		Board:  NewBoardService(repos.Board, repos.User),
		Column: NewColumnService(repos.Column, repos.Board),
		Card:   NewCardService(repos.Card, repos.Column, repos.User),
		Label:  NewLabelService(repos.Label, repos.Board),
		// Initialize other services here
	}
}