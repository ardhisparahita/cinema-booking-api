package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type GenreService interface {
	CreateGenre(ctx context.Context, req request.CreateGenreRequest) (*response.GenreResponse, error)
	GetAllGenres(ctx context.Context) ([]response.GenreResponse, error)
	GetGenreByID(ctx context.Context, id uint) (*response.GenreResponse, error)
	UpdateGenre(ctx context.Context, id uint, req request.UpdateGenreRequest) (*response.GenreResponse, error)
	DeleteGenre(ctx context.Context, id uint) error
}
