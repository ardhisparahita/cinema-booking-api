package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrRefreshNotFound    = errors.New("refresh token not found")
	ErrRefreshExpired     = errors.New("refresh token expired")
	ErrRefreshRevoked     = errors.New("refresh token revoked")
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) FindUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	err := r.DB.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return r.DB.WithContext(ctx).Create(token).Error
}

func (r *UserRepositoryImpl) FindRefreshToken(ctx context.Context, TokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken

	err := r.DB.WithContext(ctx).Where("token = ?", TokenHash).First(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRefreshNotFound
	}
	if err != nil {
		return nil, err
	}
	if token.RevokedAt != nil {
		return nil, ErrRefreshRevoked
	}
	if !token.ExpiresAt.After(time.Now()) {
		return nil, ErrRefreshExpired
	}

	return &token, nil
}

func (r *UserRepositoryImpl) RevokeRefreshToken(ctx context.Context, TokenHash string, revokedAt time.Time) error {
	result := r.DB.WithContext(ctx).Model(&models.RefreshToken{}).Where("token_hash = ? AND revoked_at IS NULL", TokenHash).Update("revoked_at", revokedAt)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRefreshNotFound
	}

	return nil
}
