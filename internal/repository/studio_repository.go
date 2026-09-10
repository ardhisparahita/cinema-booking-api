package repository

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
)

type StudioRepository interface {
	CreateStudio(ctx context.Context, studio *models.Studio) error
	FindAllStudios(ctx context.Context, theaterID uint) ([]models.Studio, error)
	FindStudioByID(ctx context.Context, id uint) (*models.Studio, error)
	FindStudioByName(ctx context.Context, theaterID uint, name string) (*models.Studio, error)
	UpdateStudio(ctx context.Context, studio *models.Studio) error
	DeleteStudio(ctx context.Context, id uint) error
}
