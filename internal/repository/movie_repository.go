package repository

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *models.Movie) error
	FindAllMovies(ctx context.Context) ([]models.Movie, error)
	FindMovieByID(ctx context.Context, id uint) (*models.Movie, error)
	UpdateMovie(ctx context.Context, movie *models.Movie) error
	DeleteMovie(ctx context.Context, id uint) error
	CreateMovieGenres(ctx context.Context, movieGenres []models.MovieGenre) error
	FindMovieGenres(ctx context.Context, movieID uint) ([]models.MovieGenre, error)
	DeleteMovieGenres(ctx context.Context, movieID uint) error
	FindGenresByIDs(ctx context.Context, genreIDs []uint) ([]models.Genre, error)
}
