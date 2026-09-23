package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	appErrors "github.com/ardhisparahita/cinema-booking-api/internal/errors"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
)

type TheaterServiceImpl struct {
	Repo repository.TheaterRepository
}

func NewTheaterService(repo repository.TheaterRepository) TheaterService {
	return &TheaterServiceImpl{
		Repo: repo,
	}
}

func (s *TheaterServiceImpl) CreateTheater(ctx context.Context, req request.CreateAndUpdateTheaterRequest) (*response.TheaterResponse, error) {
	theater := &models.Theater{
		Name:    strings.TrimSpace(req.Name),
		City:    strings.TrimSpace(req.City),
		Address: strings.TrimSpace(req.Address),
	}

	if err := s.Repo.CreateTheater(ctx, theater); err != nil {
		return nil, err
	}

	return toTheaterResponse(theater), nil
}

func (s *TheaterServiceImpl) GetAllTheaters(ctx context.Context, page, limit int) ([]response.TheaterResponse, int64, error) {
	theaters, total, err := s.Repo.FindAllTheater(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]response.TheaterResponse, 0, len(theaters))

	for _, theater := range theaters {
		result = append(result, *toTheaterResponse(&theater))
	}

	return result, total, nil
}

func (s *TheaterServiceImpl) GetTheaterByID(ctx context.Context, id uint) (*response.TheaterResponse, error) {
	theater, err := s.Repo.FindTheaterByID(ctx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrTheaterNotFound) {
			return nil, appErrors.ErrTheaterNotFound
		}
		return nil, err
	}

	return toTheaterResponse(theater), nil
}

func (s *TheaterServiceImpl) UpdateTheater(ctx context.Context, id uint, req request.CreateAndUpdateTheaterRequest) (*response.TheaterResponse, error) {
	theater, err := s.Repo.FindTheaterByID(ctx, id)
	if err != nil {
		if errors.Is(err, appErrors.ErrTheaterNotFound) {
			return nil, appErrors.ErrTheaterNotFound
		}
		return nil, err
	}

	theater.Name = strings.TrimSpace(req.Name)
	theater.City = strings.TrimSpace(req.City)
	theater.Address = strings.TrimSpace(req.Address)

	if err := s.Repo.UpdateTheater(ctx, theater); err != nil {
		return nil, err
	}

	return toTheaterResponse(theater), nil
}

func (s *TheaterServiceImpl) DeleteTheater(ctx context.Context, id uint) error {
	if _, err := s.Repo.FindTheaterByID(ctx, id); err != nil {
		if errors.Is(err, appErrors.ErrTheaterNotFound) {
			return appErrors.ErrTheaterNotFound
		}
		return err
	}

	return s.Repo.DeleteTheater(ctx, id)
}

func toTheaterResponse(theater *models.Theater) *response.TheaterResponse {
	return &response.TheaterResponse{
		ID:        theater.ID,
		Name:      theater.Name,
		City:      theater.City,
		Address:   theater.Address,
		CreatedAt: theater.CreatedAt,
		UpdatedAt: theater.UpdatedAt,
	}
}
