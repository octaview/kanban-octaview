package repository

import (
	"context"
	"errors"

	"github.com/octaview/kanban-octaview/internal/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return models.NewDatabaseError("creating user", result.Error)
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.NewDatabaseError("getting user by ID", result.Error)
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.NewDatabaseError("getting user by email", result.Error)
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return models.NewDatabaseError("updating user", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return models.NewDatabaseError("deleting user", result.Error)
	}
	if result.RowsAffected == 0 {
		return models.ErrUserNotFound
	}
	return nil
}