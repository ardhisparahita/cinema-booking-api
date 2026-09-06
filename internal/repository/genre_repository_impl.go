package repository

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrGenreNotFound     = errors.New("genre not found")
	ErrGenreAlreadyExist = errors.New("genre already exist")
)

type GenreRepositoryImpl struct {
	DB *gorm.DB
}

func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &GenreRepositoryImpl{
		DB: db,
	}
}

func (r *GenreRepositoryImpl) CreateGenre(ctx context.Context, genre *models.Genre) error {
	return r.DB.WithContext(ctx).Create(genre).Error
}

func (r *GenreRepositoryImpl) FindAllGenre(ctx context.Context) ([]models.Genre, error) {
	var genres []models.Genre

	err := r.DB.WithContext(ctx).Order("name ASD").Find(&genres).Error

	return genres, err
}

func (r *GenreRepositoryImpl) FindGenreByID(ctx context.Context, id uint) (*models.Genre, error) {
	var genre models.Genre

	err := r.DB.WithContext(ctx).First(&genre, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGenreNotFound
	}

	if err != nil {
		return nil, err
	}

	return &genre, nil
}

func (r *GenreRepositoryImpl) FindGenreByName(ctx context.Context, name string) (*models.Genre, error) {
	var genre models.Genre

	err := r.DB.WithContext(ctx).Where("name = ?", name).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGenreNotFound
	}

	if err != nil {
		return nil, err
	}

	return &genre, nil
}

func (r *GenreRepositoryImpl) UpdateGenre(ctx context.Context, genre *models.Genre) error {
	return r.DB.WithContext(ctx).Model(&models.Genre{}).Where("id = ?", genre.ID).Update("name", genre.Name).Error
}

func (r *GenreRepositoryImpl) DeleteGenre(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.Genre{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrGenreNotFound
	}

	return nil
}
