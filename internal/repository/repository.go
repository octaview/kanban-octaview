package repository

import (
	"context"

	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

type BoardRepository interface {
	Create(ctx context.Context, board *models.Board) error
	GetByID(ctx context.Context, id uint) (*models.Board, error)
	GetByOwnerID(ctx context.Context, ownerID uint) ([]models.Board, error)
	Update(ctx context.Context, board *models.Board) error
	Delete(ctx context.Context, id uint) error
}

type ColumnRepository interface {
	Create(ctx context.Context, column *models.Column) error
	GetByID(ctx context.Context, id uint) (*models.Column, error)
	GetByBoardID(ctx context.Context, boardID uint) ([]models.Column, error)
	Update(ctx context.Context, column *models.Column) error
	Delete(ctx context.Context, id uint) error
	UpdatePositions(ctx context.Context, columns []models.Column) error
}

type CardRepository interface {
	Create(ctx context.Context, card *models.Card) error
	GetByID(ctx context.Context, id uint) (*models.Card, error)
	GetByColumnID(ctx context.Context, columnID uint) ([]models.Card, error)
	Update(ctx context.Context, card *models.Card) error
	Delete(ctx context.Context, id uint) error
	UpdatePositions(ctx context.Context, cards []models.Card) error
	MoveToColumn(ctx context.Context, cardID, columnID uint, position int) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id uint) (*models.Comment, error)
	GetByCardID(ctx context.Context, cardID uint) ([]models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id uint) error
}

type LabelRepository interface {
	Create(ctx context.Context, label *models.Label) error
	GetByID(ctx context.Context, id uint) (*models.Label, error)
	GetByBoardID(ctx context.Context, boardID uint) ([]models.Label, error)
	Update(ctx context.Context, label *models.Label) error
	Delete(ctx context.Context, id uint) error
}

type CardLabelRepository interface {
	AddLabelToCard(ctx context.Context, cardID uint, labelID uint) error
	RemoveLabelFromCard(ctx context.Context, cardID uint, labelID uint) error
	GetLabelsByCardID(ctx context.Context, cardID uint) ([]models.Label, error)
	GetCardsByLabelID(ctx context.Context, labelID uint) ([]models.Card, error)
}

type Repositories struct {
	User      UserRepository
	Board     BoardRepository
	Column    ColumnRepository
	Card      CardRepository
	Comment   CommentRepository
	Label     LabelRepository
	CardLabel CardLabelRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:      NewUserRepo(db),
		Board:     NewBoardRepo(db),
		Column:    NewColumnRepo(db),
		Card:      NewCardRepo(db),
		Comment:   NewCommentRepo(db),
		Label:     NewLabelRepo(db),
		CardLabel: NewCardLabelRepo(db),
	}
}