package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"github.com/ardhisparahita/cinema-booking-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrMovieNotFound      = errors.New("movie not found")
	ErrGenreNotFound      = errors.New("one or more genres not found")
	ErrInvalidMovieRating = errors.New("invalid movie rating")
)

type MovieServiceImpl struct {
	Repo repository.MovieRepository
	DB   *gorm.DB
}

func NewMovieService(repo repository.MovieRepository, db *gorm.DB) MovieService {
	return &MovieServiceImpl{
		Repo: repo,
		DB:   db,
	}
}

func (s *MovieServiceImpl) CreateMovie(ctx context.Context, req request.CreateMovieRequest) (*response.MovieResponse, error) {
	title := strings.TrimSpace(req.Title)

	if title == "" {
		return nil, errors.New("movie title is required")
	}

	if !isValidMovieRating(req.Rating) {
		return nil, ErrInvalidMovieRating
	}

	genres, err := s.Repo.FindGenresByIDs(ctx, req.GenreIDs)
	if err != nil {
		return nil, err
	}

	if len(genres) != len(uniqueUint(req.GenreIDs)) {
		return nil, ErrGenreNotFound
	}

	movie := &models.Movie{
		Title:       title,
		Synopsis:    strings.TrimSpace(req.Synopsis),
		DurationMin: req.DurationMin,
		Rating:      req.Rating,
		PosterURL:   strings.TrimSpace(req.PosterURL),
		ReleaseDate: req.ReleaseDate,
		IsActive:    true,
	}

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewMovieRepository(tx)

		if err := txRepo.CreateMovie(ctx, movie); err != nil {
			return err
		}

		movieGenres := make([]models.MovieGenre, 0, len(req.GenreIDs))

		for _, genreID := range uniqueUint(req.GenreIDs) {
			movieGenres = append(movieGenres, models.MovieGenre{
				MovieID: movie.ID,
				GenreID: genreID,
			})
		}

		if err := txRepo.CreateMovieGenres(ctx, movieGenres); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetMovieByID(ctx, movie.ID)
}

func (s *MovieServiceImpl) GetAllMovies(ctx context.Context) ([]response.MovieResponse, error) {
	movies, err := s.Repo.FindAllMovies(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]response.MovieResponse, 0, len(movies))

	for _, movie := range movies {
		result = append(result, toMovieResponse(movie))
	}

	return result, nil
}

func (s *MovieServiceImpl) GetMovieByID(ctx context.Context, id uint) (*response.MovieResponse, error) {
	movie, err := s.Repo.FindMovieByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrMovieNotFound) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}

	res := toMovieResponse(*movie)

	return &res, nil
}

func (s *MovieServiceImpl) UpdateMovie(ctx context.Context, id uint, req request.UpdateMovieRequest) (*response.MovieResponse, error) {
	movie, err := s.Repo.FindMovieByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrMovieNotFound) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}

	if !isValidMovieRating(req.Rating) {
		return nil, ErrInvalidMovieRating
	}

	genres, err := s.Repo.FindGenresByIDs(ctx, req.GenreIDs)
	if err != nil {
		return nil, err
	}

	uniqueGenreIDs := uniqueUint(req.GenreIDs)

	if len(genres) != len(uniqueGenreIDs) {
		return nil, ErrGenreNotFound
	}

	movie.Title = strings.TrimSpace(req.Title)
	movie.Synopsis = strings.TrimSpace(req.Synopsis)
	movie.DurationMin = req.DurationMin
	movie.Rating = req.Rating
	movie.PosterURL = strings.TrimSpace(req.PosterURL)
	movie.ReleaseDate = req.ReleaseDate
	movie.IsActive = req.IsActive

	err = s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := repository.NewMovieRepository(tx)

		if err := txRepo.UpdateMovie(ctx, movie); err != nil {
			return err
		}

		if err := txRepo.DeleteMovieGenres(ctx, movie.ID); err != nil {
			return err
		}

		movieGenres := make([]models.MovieGenre, 0, len(uniqueGenreIDs))

		for _, genreID := range uniqueGenreIDs {
			movieGenres = append(movieGenres, models.MovieGenre{
				MovieID: movie.ID,
				GenreID: genreID,
			})
		}

		if err := txRepo.CreateMovieGenres(ctx, movieGenres); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetMovieByID(ctx, movie.ID)
}

func (s *MovieServiceImpl) DeleteMovie(ctx context.Context, id uint) error {
	if _, err := s.Repo.FindMovieByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrMovieNotFound) {
			return ErrMovieNotFound
		}
		return err
	}

	return s.Repo.DeleteMovie(ctx, id)
}

func isValidMovieRating(rating string) bool {
	switch rating {
	case "SU", "13+", "17+", "21+":
		return true
	default:
		return false
	}
}

func uniqueUint(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	result := make([]uint, 0, len(ids))

	for _, id := range ids {
		if _, exist := seen[id]; exist {
			continue
		}

		seen[id] = struct{}{}
		result = append(result, id)
	}

	return result
}

func toMovieResponse(movie models.Movie) response.MovieResponse {
	genres := make([]response.GenreResponse, 0, len(movie.MovieGenres))

	for _, movieGenre := range movie.MovieGenres {
		genres = append(genres, response.GenreResponse{
			ID:   movieGenre.Genre.ID,
			Name: movieGenre.Genre.Name,
		})
	}

	return response.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Synopsis:    movie.Synopsis,
		DurationMin: movie.DurationMin,
		Rating:      movie.Rating,
		PosterURL:   movie.PosterURL,
		ReleaseDate: movie.ReleaseDate,
		IsActive:    movie.IsActive,
		Genres:      genres,
		CreatedAt:   movie.CreatedAt,
		UpdatedAt:   movie.UpdatedAt,
	}
}
