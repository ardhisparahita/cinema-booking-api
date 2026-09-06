package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
)

var (
	ErrGenreNotFound     = errors.New("genre not found")
	ErrGenreAlreadyExist = errors.New("genre already exist")
)

type GenreServiceImpl struct {
	Repo repository.GenreRepository
}

func NewGenreService(repo repository.GenreRepository) GenreService {
	return &GenreServiceImpl{
		Repo: repo,
	}
}

func (s *GenreServiceImpl) CreateGenre(ctx context.Context, req request.CreateGenreRequest) (*response.GenreResponse, error) {
	name := strings.TrimSpace(req.Name)

	_, err := s.Repo.FindGenreByName(ctx, name)
	if err == nil {
		return nil, ErrGenreAlreadyExist
	}

	if !errors.Is(err, repository.ErrGenreNotFound) {
		return nil, err
	}

	genre := &models.Genre{
		Name: name,
	}

	if err := s.Repo.CreateGenre(ctx, genre); err != nil {
		return nil, err
	}

	return &response.GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}, nil
}

func (s *GenreServiceImpl) GetAllGenres(ctx context.Context) ([]response.GenreResponse, error) {
	genres, err := s.Repo.FindAllGenre(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]response.GenreResponse, 0, len(genres))

	for _, genre := range genres {
		result = append(result, response.GenreResponse{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	return result, nil
}

func (s *GenreServiceImpl) GetGenreByID(ctx context.Context, id uint) (*response.GenreResponse, error) {
	genre, err := s.Repo.FindGenreByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrGenreNotFound) {
			return nil, ErrGenreNotFound
		}
		return nil, err
	}

	return &response.GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}, nil
}

func (s *GenreServiceImpl) UpdateGenre(ctx context.Context, id uint, req request.UpdateGenreRequest) (*response.GenreResponse, error) {
	name := strings.TrimSpace(req.Name)

	genre, err := s.Repo.FindGenreByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrGenreNotFound) {
			return nil, ErrGenreNotFound
		}
	}
	return nil, err

	existing, err := s.Repo.FindGenreByName(ctx, name)
	if err == nil && existing.ID != genre.ID {
		return nil, ErrGenreAlreadyExist
	}

	if err != nil && !errors.Is(err, repository.ErrGenreNotFound) {
		return nil, err
	}

	genre.Name = name

	if err := s.Repo.UpdateGenre(ctx, genre); err != nil {
		return nil, err
	}

	return &response.GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}, nil
}

func (s *GenreServiceImpl) DeleteGenre(ctx context.Context, id uint) error {
	if _, err := s.Repo.FindGenreByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrGenreNotFound) {
			return ErrGenreNotFound
		}

		return err
	}

	return s.Repo.DeleteGenre(ctx, id)
}
