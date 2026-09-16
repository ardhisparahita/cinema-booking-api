package repository

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrSeatNotFound     = errors.New("seat not found")
	ErrSeatAlreadyExist = errors.New("seat already exist")
)

type SeatRepositoryImpl struct {
	DB *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &SeatRepositoryImpl{
		DB: db,
	}
}

func (r *SeatRepositoryImpl) CreateSeat(ctx context.Context, seat *models.Seat) error {
	var existing models.Seat

	err := r.DB.WithContext(ctx).Where("studio_id = ? AND row_label = ? AND col_number = ?", seat.StudioID, seat.RowLabel, seat.ColNumber).First(&existing).Error

	if err == nil {
		return ErrSeatAlreadyExist
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.DB.WithContext(ctx).Create(seat).Error
}

func (r *SeatRepositoryImpl) FindAllSeats(ctx context.Context, studioID uint) ([]models.Seat, error) {
	var seats []models.Seat

	err := r.DB.WithContext(ctx).Where("studio_id = ?", studioID).Order("row_label ASC, col_number ASC").Find(&seats).Error

	return seats, err
}

func (r *SeatRepositoryImpl) FindSeatByID(ctx context.Context, id uint) (*models.Seat, error) {
	var seat models.Seat

	err := r.DB.WithContext(ctx).Find(&seat, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSeatNotFound
	}

	if err != nil {
		return nil, err
	}

	return &seat, nil
}

func (r *SeatRepositoryImpl) FindSeatByPosition(ctx context.Context, studioID uint, rowLabel string, colNumber uint16) (*models.Seat, error) {
	var seat models.Seat

	err := r.DB.WithContext(ctx).Where("studio_id = ? AND row_label = ? AND col_number = ?", studioID, rowLabel, colNumber).First(&seat).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSeatNotFound
	}

	if err != nil {
		return nil, err
	}

	return &seat, nil
}

func (r *SeatRepositoryImpl) UpdateSeat(ctx context.Context, seat *models.Seat) error {
	return r.DB.WithContext(ctx).Model(&models.Seat{}).Where("id = ?", seat.ID).Updates(map[string]any{
		"row_label":  seat.RowLabel,
		"col_number": seat.ColNumber,
		"seat_type":  seat.SeatType,
	}).Error
}

func (r *SeatRepositoryImpl) DeleteSeat(ctx context.Context, id uint16) error {
	result := r.DB.WithContext(ctx).Delete(&models.Seat{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrSeatNotFound
	}

	return nil
}

func (r *SeatRepositoryImpl) FindSeatByIDs(ctx context.Context, ids []uint) ([]models.Seat, error) {
	var seats []models.Seat

	if len(seats) == 0 {
		return seats, nil
	}

	err := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&seats).Error

	return seats, err
}
