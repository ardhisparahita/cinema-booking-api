package service

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
)

var (
	ErrShowtimeNotFound       = errors.New("showtime not found")
	ErrMovieNotFoundShowtime  = errors.New("movie not found")
	ErrStudioNotFoundShowtime = errors.New("studio not found")
	ErrInvalidShowtimeTime    = errors.New("end time must be after start time")
	ErrShowtimeConflict       = errors.New("showtime conflict with another showtime")
	ErrInvalidShowtimePrice   = errors.New("price must be greater than zero")
)

type ShowtimeServiceImpl struct {
	Repo       repository.ShowtimeRepository
	MovieRepo  repository.MovieRepository
	StudioRepo repository.StudioRepository
}

func NewShowtimeService(repo repository.ShowtimeRepository, movieRepo repository.MovieRepository, studioRepo repository.StudioRepository) ShowtimeService {
	return &ShowtimeServiceImpl{
		Repo:       repo,
		MovieRepo:  movieRepo,
		StudioRepo: studioRepo,
	}
}

func (s *ShowtimeServiceImpl) CreateShowtime(ctx context.Context, req request.CreateAndUpdateShowtimeRequest) (*response.ShowtimeResponse, error) {
	if !req.EndTime.After(req.StartTime) {
		return nil, ErrInvalidShowtimeTime
	}

	if req.Price <= 0 {
		return nil, ErrInvalidShowtimePrice
	}

	if _, err := s.MovieRepo.FindMovieByID(ctx, req.MovieID); err != nil {
		if errors.Is(err, repository.ErrMovieNotFound) {
			return nil, ErrMovieNotFoundShowtime
		}

		return nil, err
	}

	if _, err := s.StudioRepo.FindStudioByID(ctx, req.StudioID); err != nil {
		if errors.Is(err, repository.ErrStudioNotFound) {
			return nil, ErrStudioNotFoundShowtime
		}
		return nil, err
	}

	conflicts, err := s.Repo.FindShowtimeByStudioAndTime(ctx, req.StudioID, req.StartTime, req.EndTime, 0)
	if err != nil {
		return nil, err
	}

	if len(conflicts) > 0 {
		return nil, ErrShowtimeConflict
	}

	showtime := &models.Showtime{
		MovieID:   req.MovieID,
		StudioID:  req.StudioID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Price:     req.Price,
	}
	if err := s.Repo.CreateShowtime(ctx, showtime); err != nil {
		return nil, err
	}

	return toShowtimeResponse(showtime), nil
}

func (s *ShowtimeServiceImpl) GetAllShowtimes(ctx context.Context) ([]response.ShowtimeResponse, error) {
	showtimes, err := s.Repo.FindAllShowtimes(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]response.ShowtimeResponse, 0, len(showtimes))

	for _, showtime := range showtimes {
		result = append(result, *toShowtimeResponse(&showtime))
	}

	return result, nil
}

func (s *ShowtimeServiceImpl) GetShowtimeByID(ctx context.Context, id uint) (*response.ShowtimeResponse, error) {
	showtime, err := s.Repo.FindShowtimeByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrShowtimeNotFound) {
			return nil, ErrShowtimeNotFound
		}
		return nil, err
	}

	return toShowtimeResponse(showtime), nil
}

func (s *ShowtimeServiceImpl) UpdateShowtime(ctx context.Context, id uint, req request.CreateAndUpdateShowtimeRequest) (*response.ShowtimeResponse, error) {
	showtime, err := s.Repo.FindShowtimeByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrShowtimeNotFound) {
			return nil, ErrShowtimeNotFound
		}
		return nil, err
	}

	if !req.StartTime.After(req.StartTime) {
		return nil, ErrInvalidShowtimeTime
	}

	if req.Price <= 0 {
		return nil, ErrInvalidShowtimePrice
	}

	if _, err := s.MovieRepo.FindMovieByID(ctx, req.MovieID); err != nil {
		if errors.Is(err, repository.ErrMovieNotFound) {
			return nil, ErrMovieNotFoundShowtime
		}

		return nil, err
	}

	if _, err := s.StudioRepo.FindStudioByID(ctx, req.StudioID); err != nil {
		if errors.Is(err, repository.ErrStudioNotFound) {
			return nil, ErrStudioNotFoundShowtime
		}
		return nil, err
	}

	conflicts, err := s.Repo.FindShowtimeByStudioAndTime(ctx, req.StudioID, req.StartTime, req.EndTime, 0)
	if err != nil {
		return nil, err
	}

	if len(conflicts) > 0 {
		return nil, ErrShowtimeConflict
	}

	showtime.MovieID = req.MovieID
	showtime.StudioID = req.StudioID
	showtime.StartTime = req.StartTime
	showtime.EndTime = req.EndTime
	showtime.Price = req.Price

	if err := s.Repo.UpdateShowtime(ctx, showtime); err != nil {
		return nil, err
	}

	return toShowtimeResponse(showtime), nil
}

func (s *ShowtimeServiceImpl) DeleteShowtime(ctx context.Context, id uint) error {
	if _, err := s.Repo.FindShowtimeByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrShowtimeNotFound) {
			return ErrShowtimeNotFound
		}
		return err
	}

	return s.Repo.DeleteShowtime(ctx, id)
}

func toShowtimeResponse(showtime *models.Showtime) *response.ShowtimeResponse {
	return &response.ShowtimeResponse{
		ID:        showtime.ID,
		MovieID:   showtime.MovieID,
		StudioID:  showtime.StudioID,
		StartTime: showtime.StartTime,
		EndTime:   showtime.EndTime,
		Price:     showtime.Price,
		CreatedAt: showtime.CreatedAt,
		UpdatedAt: showtime.UpdatedAt,
	}
}
