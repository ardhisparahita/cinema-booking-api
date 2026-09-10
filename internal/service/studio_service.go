package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type StudioService interface {
	CreateStudio(ctx context.Context, theaterID uint, req request.CreateAndUpdateStudioRequest) (*response.StudioResponse, error)
	GetAllStudios(ctx context.Context, theaterID uint) ([]response.StudioResponse, error)
	GetStudioByID(ctx context.Context, id uint) (*response.StudioResponse, error)
	UpdateStudio(ctx context.Context, id uint, req request.CreateAndUpdateStudioRequest) (*response.StudioResponse, error)
	DeleteStudio(ctx context.Context, id uint) error
}
