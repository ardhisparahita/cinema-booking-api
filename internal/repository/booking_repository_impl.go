package repository

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrBookingNotFound = errors.New("booking not found")
)

type BookingRepositoryImpl struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &BookingRepositoryImpl{
		DB: db,
	}
}

func (r *BookingRepositoryImpl) CreateBooking(ctx context.Context, tx *gorm.DB, booking *models.Booking) error {
	return tx.WithContext(ctx).Create(booking).Error
}

func (r *BookingRepositoryImpl) CreateBookingSeats(ctx context.Context, tx *gorm.DB, bookingSeats []models.BookingSeat) error {
	if len(bookingSeats) == 0 {
		return nil
	}

	return tx.WithContext(ctx).Create(&bookingSeats).Error
}

func (r *BookingRepositoryImpl) FindBookingByID(ctx context.Context, id uint) (*models.Booking, error) {
	var booking models.Booking

	err := r.DB.WithContext(ctx).First(&booking, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrBookingNotFound
	}

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepositoryImpl) FindBookingByUserID(ctx context.Context, userID uint) ([]models.Booking, error) {
	var bookings []models.Booking

	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&bookings).Error

	return bookings, err
}

func (r *BookingRepositoryImpl) FindBookedSeatIDs(ctx context.Context, showtimeID uint, seatIDs []uint) ([]uint, error) {
	if len(seatIDs) == 0 {
		return []uint{}, nil
	}

	var bookedSeatIDs []uint

	err := r.DB.WithContext(ctx).Model(&models.BookingSeat{}).Where("showtime_id = ?", showtimeID).Where("seat_id = ?", seatIDs).Pluck("seat_id", &bookedSeatIDs).Error

	return bookedSeatIDs, err
}

func (r *BookingRepositoryImpl) CancelBooking(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Model(&models.Booking{}).Where("id = ?", id).Update("status", "cancelled")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrBookingNotFound
	}

	return nil
}
