package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type SeatService interface {
	CreateSeat(ctx context.Context, studioID uint, req request.CreateAndUpdateSeatRequest) (*response.SeatResponse, error)
	GetAllSeats(ctx context.Context, studioID uint) ([]response.SeatResponse, error)
	GetSetByID(ctx context.Context, id uint) (*response.SeatResponse, error)
	UpdateSeat(ctx context.Context, id uint, req request.CreateAndUpdateSeatRequest) (*response.SeatResponse, error)
	DeleteSeat(ctx context.Context, id uint) error
}
