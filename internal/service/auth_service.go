package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type AuthService interface {
	Register(ctx context.Context, req request.RegisterRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error)
	Refresh(ctx context.Context, rawRefreshToken string) (*response.TokenResponse, error)
	Logout(ctx context.Context, rawRefreshToken string) error
}
