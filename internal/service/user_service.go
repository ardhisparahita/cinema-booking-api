package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*response.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req request.UpdateUserRequest) (*response.UserResponse, error)
	DeleteProfile(ctx context.Context, userID uint) error
}
