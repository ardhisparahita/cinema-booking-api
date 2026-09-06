package repository

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
)

type GenreRepository interface {
	CreateGenre(ctx context.Context, genre *models.Genre) error
	FindAllGenre(ctx context.Context) ([]models.Genre, error)
	FindGenreByID(ctx context.Context, id uint) (*models.Genre, error)
	FindGenreByName(ctx context.Context, name string) (*models.Genre, error)
	UpdateGenre(ctx context.Context, genre *models.Genre) error
	DeleteGenre(ctx context.Context, id uint) error
}
