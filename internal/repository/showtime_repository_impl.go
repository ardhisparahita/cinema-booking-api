package repository

import (
	"context"
	"errors"
	"time"

	appErrors "github.com/ardhisparahita/cinema-booking-api/internal/errors"
	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

type ShowtimeRepositoryImpl struct {
	DB *gorm.DB
}

func NewShowtimeRepository(db *gorm.DB) ShowtimeRepository {
	return &ShowtimeRepositoryImpl{
		DB: db,
	}
}

func (r *ShowtimeRepositoryImpl) CreateShowtime(ctx context.Context, showtime *models.Showtime) error {
	return r.DB.WithContext(ctx).Create(showtime).Error
}

func (r *ShowtimeRepositoryImpl) FindAllShowtimes(ctx context.Context, page, limit int) ([]models.Showtime, int64, error) {
	var showtimes []models.Showtime
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.Showtime{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := r.DB.WithContext(ctx).Order("start_time ASC").Limit(limit).Offset(offset).Find(&showtimes).Error

	return showtimes, total, err
}

func (r *ShowtimeRepositoryImpl) FindShowtimeByID(ctx context.Context, id uint) (*models.Showtime, error) {
	var showtime models.Showtime

	err := r.DB.WithContext(ctx).First(&showtime, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.ErrShowtimeNotFound
	}

	if err != nil {
		return nil, err
	}

	return &showtime, nil
}

func (r *ShowtimeRepositoryImpl) FindShowtimeByStudioAndTime(ctx context.Context, studioID uint, startTime time.Time, endTime time.Time, excludeID uint) ([]models.Showtime, error) {
	var showtimes []models.Showtime

	query := r.DB.WithContext(ctx).Where("studio_id = ?", studioID).Where("start_time < ?", endTime).Where("end_time > ?", startTime)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Find(&showtimes).Error

	return showtimes, err
}

func (r *ShowtimeRepositoryImpl) UpdateShowtime(ctx context.Context, showtime *models.Showtime) error {
	return r.DB.WithContext(ctx).Model(&models.Showtime{}).Where("id = ?", showtime.ID).Updates(map[string]any{
		"movie_id":   showtime.MovieID,
		"studio_id":  showtime.StudioID,
		"start_time": showtime.StartTime,
		"end_time":   showtime.EndTime,
		"price":      showtime.Price,
	}).Error
}

func (r *ShowtimeRepositoryImpl) DeleteShowtime(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.Showtime{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return appErrors.ErrShowtimeNotFound
	}

	return nil
}
