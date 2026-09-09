package repository

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
)

type TheaterRepository interface {
	CreateTheater(ctx context.Context, theater *models.Theater) error
	FindAllTheater(ctx context.Context) ([]models.Theater, error)
	FindTheaterByID(ctx context.Context, id uint) (*models.Theater, error)
	UpdateTheater(ctx context.Context, theater *models.Theater) error
	DeleteTheater(ctx context.Context, id uint) error
}
