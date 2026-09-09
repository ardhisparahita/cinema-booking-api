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
	ErrTheaterNotFound = errors.New("theater not found")
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

func (s *TheaterServiceImpl) GetAllTheaters(ctx context.Context) ([]response.TheaterResponse, error) {
	theaters, err := s.Repo.FindAllTheater(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]response.TheaterResponse, 0, len(theaters))

	for _, theater := range theaters {
		result = append(result, *toTheaterResponse(&theater))
	}

	return result, nil
}

func (s *TheaterServiceImpl) GetTheaterByID(ctx context.Context, id uint) (*response.TheaterResponse, error) {
	theater, err := s.Repo.FindTheaterByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTheaterNotFound) {
			return nil, ErrTheaterNotFound
		}
		return nil, err
	}

	return toTheaterResponse(theater), nil
}

func (s *TheaterServiceImpl) UpdateTheater(ctx context.Context, id uint, req request.CreateAndUpdateTheaterRequest) (*response.TheaterResponse, error) {
	theater, err := s.Repo.FindTheaterByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTheaterNotFound) {
			return nil, ErrTheaterNotFound
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
		if errors.Is(err, repository.ErrTheaterNotFound) {
			return ErrTheaterNotFound
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
