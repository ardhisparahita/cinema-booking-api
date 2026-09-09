package repository

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrTheaterNotFound = errors.New("theater not found")
)

type TheaterRepositoryImpl struct {
	DB *gorm.DB
}

func NewTheaterRepository(db *gorm.DB) TheaterRepository {
	return &TheaterRepositoryImpl{
		DB: db,
	}
}

func (r *TheaterRepositoryImpl) CreateTheater(ctx context.Context, theater *models.Theater) error {
	return r.DB.WithContext(ctx).Create(theater).Error
}

func (r *TheaterRepositoryImpl) FindAllTheater(ctx context.Context) ([]models.Theater, error) {
	var theaters []models.Theater

	err := r.DB.WithContext(ctx).Order("name ASC").Find(&theaters).Error

	return theaters, err
}

func (r *TheaterRepositoryImpl) FindTheaterByID(ctx context.Context, id uint) (*models.Theater, error) {
	var theater models.Theater

	err := r.DB.WithContext(ctx).First(&theater, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTheaterNotFound
	}

	if err != nil {
		return nil, err
	}

	return &theater, nil
}

func (r *TheaterRepositoryImpl) UpdateTheater(ctx context.Context, theater *models.Theater) error {
	return r.DB.WithContext(ctx).Model(&models.Theater{}).Where("id = ?", theater.ID).Updates(map[string]any{
		"name":    theater.Name,
		"city":    theater.City,
		"address": theater.Address,
	}).Error
}
func (r *TheaterRepositoryImpl) DeleteTheater(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.Theater{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTheaterNotFound
	}

	return nil
}
