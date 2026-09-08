package repository

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrMovieNotFound = errors.New("movie not found")
)

type MovieRepositoryImpl struct {
	DB *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &MovieRepositoryImpl{
		DB: db,
	}
}

func (r *MovieRepositoryImpl) CreateMovie(ctx context.Context, movie *models.Movie) error {
	return r.DB.WithContext(ctx).Create(movie).Error
}

func (r *MovieRepositoryImpl) FindAllMovies(ctx context.Context) ([]models.Movie, error) {
	var movies []models.Movie

	err := r.DB.WithContext(ctx).Preload("MovieGenres").Preload("MovieGenres.Genre").Order("created_at DESC").Find(&movies).Error

	return movies, err
}
func (r *MovieRepositoryImpl) FindMovieByID(ctx context.Context, id uint) (*models.Movie, error) {
	var movie models.Movie

	err := r.DB.WithContext(ctx).Preload("MovieGenres").Preload("MovieGenres.Genre").First(&movie, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGenreNotFound
	}

	if err != nil {
		return nil, err
	}

	return &movie, err
}

func (r *MovieRepositoryImpl) UpdateMovie(ctx context.Context, movie *models.Movie) error {
	return r.DB.WithContext(ctx).Model(&models.Movie{}).Where("id = ?", movie.ID).Updates(map[string]any{
		"title":        movie.Title,
		"synopsis":     movie.Synopsis,
		"duration_min": movie.DurationMin,
		"rating":       movie.Rating,
		"poster_url":   movie.PosterUrl,
		"release_date": movie.ReleaseDate,
		"is_active":    movie.IsActive,
	}).Error
}
func (r *MovieRepositoryImpl) DeleteMovie(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.Movie{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrMovieNotFound
	}

	return nil
}

func (r *MovieRepositoryImpl) CreateMovieGenres(ctx context.Context, movieGenres []models.MovieGenre) error {
	if len(movieGenres) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&movieGenres).Error
}

func (r *MovieRepositoryImpl) FindMovieGenres(ctx context.Context, movieID uint) ([]models.MovieGenre, error) {
	var movieGenres []models.MovieGenre

	err := r.DB.WithContext(ctx).Preload("Genre").Where("movie_id = ?", movieID).Find(&movieGenres).Error

	return movieGenres, err
}

func (r *MovieRepositoryImpl) DeleteMovieGenres(ctx context.Context, movieID uint) error {
	return r.DB.WithContext(ctx).Where("movie_id = ?", movieID).Delete(&models.MovieGenre{}).Error
}

func (r *MovieRepositoryImpl) FindGenresByIDs(ctx context.Context, genreIDs []uint) ([]models.Genre, error) {
	var genres []models.Genre

	if len(genreIDs) == 0 {
		return genres, nil
	}

	err := r.DB.WithContext(ctx).Where("id IN ?", genreIDs).Find(&genres).Error

	return genres, err
}
