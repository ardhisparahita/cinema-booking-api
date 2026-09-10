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
	ErrStudioNotFound           = errors.New("studio not found")
	ErrStudioAlreadyExist       = errors.New("studio already exist")
	ErrTheaterNotFoundForStudio = errors.New("theater not found")
)

type StudioServiceImpl struct {
	Repo        repository.StudioRepository
	TheaterRepo repository.TheaterRepository
}

func NewStudioService(repo repository.StudioRepository, theaterRepo repository.TheaterRepository) StudioService {
	return &StudioServiceImpl{
		Repo:        repo,
		TheaterRepo: theaterRepo,
	}
}

func (s *StudioServiceImpl) CreateStudio(ctx context.Context, theaterID uint, req request.CreateAndUpdateStudioRequest) (*response.StudioResponse, error) {
	if _, err := s.TheaterRepo.FindTheaterByID(ctx, theaterID); err != nil {
		if errors.Is(err, repository.ErrTheaterNotFound) {
			return nil, ErrTheaterNotFoundForStudio
		}
		return nil, err
	}

	name := strings.TrimSpace(req.Name)

	existing, err := s.Repo.FindStudioByName(ctx, theaterID, name)
	if err == nil && existing != nil {
		return nil, ErrStudioAlreadyExist
	}

	if err != nil && !errors.Is(err, repository.ErrStudioNotFound) {
		return nil, err
	}

	studio := &models.Studio{
		TheaterID: theaterID,
		Name:      name,
		TotalRows: req.TotalRows,
		TotalCols: req.TotalCols,
	}

	if err := s.Repo.CreateStudio(ctx, studio); err != nil {
		return nil, err
	}

	return toStudioResponse(studio), nil
}

func (s *StudioServiceImpl) GetAllStudios(ctx context.Context, theaterID uint) ([]response.StudioResponse, error) {
	if _, err := s.TheaterRepo.FindTheaterByID(ctx, theaterID); err != nil {
		if errors.Is(err, repository.ErrStudioNotFound) {
			return nil, ErrTheaterNotFoundForStudio
		}

		return nil, err
	}

	studios, err := s.Repo.FindAllStudios(ctx, theaterID)
	if err != nil {
		return nil, err
	}

	result := make([]response.StudioResponse, 0, len(studios))

	for _, studio := range studios {
		result = append(result, *toStudioResponse(&studio))
	}

	return result, nil
}

func (s *StudioServiceImpl) GetStudioByID(ctx context.Context, id uint) (*response.StudioResponse, error) {
	studio, err := s.Repo.FindStudioByID(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrStudioNotFound) {
			return nil, ErrStudioNotFound
		}
		return nil, err
	}

	return toStudioResponse(studio), nil
}

func (s *StudioServiceImpl) UpdateStudio(ctx context.Context, id uint, req request.CreateAndUpdateStudioRequest) (*response.StudioResponse, error) {
	studio, err := s.Repo.FindStudioByID(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrStudioNotFound) {
			return nil, ErrStudioNotFound
		}
		return nil, err
	}

	name := strings.TrimSpace(req.Name)

	existing, err := s.Repo.FindStudioByName(ctx, studio.TheaterID, name)
	if err == nil && existing.ID != studio.ID {
		return nil, ErrStudioAlreadyExist
	}

	if err != nil && !errors.Is(err, repository.ErrStudioNotFound) {
		return nil, err
	}

	studio.Name = name
	studio.TotalRows = req.TotalRows
	studio.TotalCols = req.TotalCols

	if err := s.Repo.UpdateStudio(ctx, studio); err != nil {
		return nil, err
	}

	return toStudioResponse(studio), nil
}

func (s *StudioServiceImpl) DeleteStudio(ctx context.Context, id uint) error {
	if _, err := s.Repo.FindStudioByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrStudioNotFound) {
			return ErrStudioNotFound
		}
		return err
	}

	return s.Repo.DeleteStudio(ctx, id)
}

func toStudioResponse(studio *models.Studio) *response.StudioResponse {
	return &response.StudioResponse{
		ID:        studio.ID,
		TheaterID: studio.TheaterID,
		Name:      studio.Name,
		TotalCols: studio.TotalCols,
		TotalRows: studio.TotalRows,
		CreatedAt: studio.CreatedAt,
		UpdatedAt: studio.UpdatedAt,
	}
}
