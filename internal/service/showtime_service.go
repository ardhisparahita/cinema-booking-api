package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type ShowtimeService interface {
	CreateShowtime(ctx context.Context, req request.CreateAndUpdateShowtimeRequest) (*response.ShowtimeResponse, error)
	GetAllShowtimes(ctx context.Context) ([]response.ShowtimeResponse, error)
	GetShowtimeByID(ctx context.Context, id uint) (*response.ShowtimeResponse, error)
	UpdateShowtime(ctx context.Context, id uint, req request.CreateAndUpdateShowtimeRequest) (*response.ShowtimeResponse, error)
	DeleteShowtime(ctx context.Context, id uint) error
}
