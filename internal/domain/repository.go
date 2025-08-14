package domain

import (
	"context"
	"time"

	"apiprofile/internal/domain/models"
)

type RefreshTokenRow struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id string) error

	CreateRefreshToken(ctx context.Context, token string, userID string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshTokenRow, error)
	UpdateRefreshToken(ctx context.Context, oldToken string, newToken string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, token string) error
}
