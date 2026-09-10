package repository

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrStudioNotFound = errors.New("studio not found")
	ErrStudioExist    = errors.New("studio already exist")
)

type StudioRepositoryImpl struct {
	DB *gorm.DB
}

func NewStudioRepository(db *gorm.DB) StudioRepository {
	return &StudioRepositoryImpl{
		DB: db,
	}
}

func (r *StudioRepositoryImpl) CreateStudio(ctx context.Context, studio *models.Studio) error {
	var existing models.Studio

	err := r.DB.WithContext(ctx).Where("theater_id = ? AND name = ?", studio.TheaterID, studio.Name).First(&existing).Error

	if err == nil {
		return ErrStudioExist
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.DB.WithContext(ctx).Create(studio).Error
}

func (r *StudioRepositoryImpl) FindAllStudios(ctx context.Context, theaterID uint) ([]models.Studio, error) {
	var studios []models.Studio

	err := r.DB.WithContext(ctx).Where("theater_id = ?", theaterID).Order("name ASC").Find(&studios).Error

	return studios, err
}

func (r *StudioRepositoryImpl) FindStudioByID(ctx context.Context, id uint) (*models.Studio, error) {
	var studio models.Studio

	err := r.DB.WithContext(ctx).First(&studio, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStudioNotFound
	}

	if err != nil {
		return nil, err
	}

	return &studio, nil
}

func (r *StudioRepositoryImpl) FindStudioByName(ctx context.Context, theaterID uint, name string) (*models.Studio, error) {
	var studio models.Studio

	err := r.DB.WithContext(ctx).Where("theater_id = ? AND name = ?", theaterID, name).First(&studio).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrStudioNotFound
	}

	if err != nil {
		return nil, err
	}

	return &studio, nil
}

func (r *StudioRepositoryImpl) UpdateStudio(ctx context.Context, studio *models.Studio) error {
	var existing models.Studio

	err := r.DB.WithContext(ctx).Where("theater_id = ? AND name = ? AND id != ?", studio.TheaterID, studio.Name, studio.ID).First(&existing).Error

	if err == nil {
		return ErrStudioExist
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	result := r.DB.WithContext(ctx).Model(&models.Studio{}).Where("id = ?", studio.ID).Updates(map[string]any{
		"name":       studio.Name,
		"total_rows": studio.TotalRows,
		"total_cols": studio.TotalCols,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrStudioNotFound
	}

	return nil
}

func (r *StudioRepositoryImpl) DeleteStudio(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.Studio{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrStudioNotFound
	}

	return nil
}
