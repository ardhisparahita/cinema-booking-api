package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	"github.com/ardhisparahita/cinema-booking-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	repo       repository.UserRepository
	jwtManager jwt.Manager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(
	repo repository.UserRepository,
	jwtManager jwt.Manager,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) AuthService {
	return &AuthServiceImpl{
		repo:       repo,
		jwtManager: jwtManager,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req request.RegisterRequest) (*response.AuthResponse, error) {
	_, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, repository.ErrEmailAlreadyExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := s.repo.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	return s.issueAuthResponse(ctx, &user)
}

func (s *AuthServiceImpl) Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return s.issueAuthResponse(ctx, user)
}

func (s *AuthServiceImpl) Refresh(ctx context.Context, rawRefreshToken string) (*response.TokenResponse, error) {
	if rawRefreshToken == "" {
		return nil, errors.New("invalid token")
	}

	tokenHash := hashToken(rawRefreshToken)
	oldToken, err := s.repo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.FindUserByID(ctx, oldToken.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Role, s.accessTTL)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if err := s.repo.RevokeRefreshToken(ctx, tokenHash, now); err != nil {
		return nil, err
	}

	if err := s.repo.CreateRefreshToken(ctx, &models.RefreshToken{
		UserID:    user.ID,
		Token:     hashToken(newRefreshToken),
		ExpiresAt: now.Add(s.refreshTTL),
	}); err != nil {
		return nil, err
	}

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
	}, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return errors.New("invalid token input")
	}

	return s.repo.RevokeRefreshToken(ctx, hashToken(rawRefreshToken), time.Now())
}

func (s *AuthServiceImpl) issueAuthResponse(ctx context.Context, user *models.User) (*response.AuthResponse, error) {
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Role, s.accessTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateRefreshToken(ctx, &models.RefreshToken{
		UserID:    user.ID,
		Token:     hashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}); err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		User: response.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Tokens: response.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(s.accessTTL.Seconds()),
		},
	}, nil
}

func generateRandomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
