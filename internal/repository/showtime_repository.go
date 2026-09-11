package repository

import (
	"context"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
)

type ShowtimeRepository interface {
	CreateShowtime(ctx context.Context, showtime *models.Showtime) error
	FindAllShowtimes(ctx context.Context) ([]models.Showtime, error)
	FindShowtimeByID(ctx context.Context, id uint) (*models.Showtime, error)
	FindShowtimeByStudioAndTime(ctx context.Context, studioID uint, startTime time.Time, endTime time.Time, excludeID uint) ([]models.Showtime, error)
	UpdateShowtime(ctx context.Context, showtime *models.Showtime) error
	DeleteShowtime(ctx context.Context, id uint) error
}
