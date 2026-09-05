package repository

import (
	"context"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, id uint) (*models.User, error)
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	FindRefreshToken(ctx context.Context, TokenHash string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, TokenHash string, revokedAt time.Time) error

	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uint) error
}
