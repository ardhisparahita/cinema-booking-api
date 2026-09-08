package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type MovieService interface {
	CreateMovie(ctx context.Context, req request.CreateMovieRequest) (*response.MovieResponse, error)
	GetAllMovies(ctx context.Context) ([]response.MovieResponse, error)
	GetMovieByID(ctx context.Context, id uint) (*response.MovieResponse, error)
	UpdateMovie(ctx context.Context, id uint, req request.UpdateMovieRequest) (*response.MovieResponse, error)
	DeleteMovie(ctx context.Context, id uint) error
}
