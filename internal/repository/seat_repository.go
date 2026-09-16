package repository

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
)

type SeatRepository interface {
	CreateSeat(ctx context.Context, seat *models.Seat) error
	FindAllSeats(ctx context.Context, studioID uint) ([]models.Seat, error)
	FindSeatByID(ctx context.Context, id uint) (*models.Seat, error)
	FindSeatByPosition(ctx context.Context, studioID uint, rowLabel string, colNumber uint16) (*models.Seat, error)
	UpdateSeat(ctx context.Context, seat *models.Seat) error
	DeleteSeat(ctx context.Context, id uint16) error
	FindSeatByIDs(ctx context.Context, ids []uint) ([]models.Seat, error)
}
