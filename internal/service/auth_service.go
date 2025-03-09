package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/octaview/kanban-octaview/internal/config"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}

type AuthService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, user *models.User) (uint, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return 0, models.ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return "", models.ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", models.ErrInvalidCredentials
	}

	token, err := s.generateJWTToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) ParseToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, models.ErrUnauthorized
	}

	return claims.UserID, nil
}

func (s *AuthService) generateJWTToken(userID uint) (string, error) {
	now := time.Now()
	claims := &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(s.cfg.JWT.ExpiresIn).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *AuthService) RefreshToken(ctx context.Context, userID uint) (string, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	token, err := s.generateJWTToken(userID)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return token, nil
}
