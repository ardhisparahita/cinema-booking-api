package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type TheaterService interface {
	CreateTheater(ctx context.Context, req request.CreateAndUpdateTheaterRequest) (*response.TheaterResponse, error)
	GetAllTheaters(ctx context.Context) ([]response.TheaterResponse, error)
	GetTheaterByID(ctx context.Context, id uint) (*response.TheaterResponse, error)
	UpdateTheater(ctx context.Context, id uint, req request.CreateAndUpdateTheaterRequest) (*response.TheaterResponse, error)
	DeleteTheater(ctx context.Context, id uint) error
}
