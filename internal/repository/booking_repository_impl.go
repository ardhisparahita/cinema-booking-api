package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrBookingNotFound   = errors.New("booking not found")
	ErrBookingNotPending = errors.New("booking is not pending")
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

	err := r.DB.WithContext(ctx).Preload("BookingSeats").Preload("BookingSeats.Seat").First(&booking, id).Error
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

	err := r.DB.WithContext(ctx).Preload("BookingSeats").Preload("BookingSeats.Seat").Where("user_id = ?", userID).Order("created_at DESC").Find(&bookings).Error

	return bookings, err
}

func (r *BookingRepositoryImpl) FindBookedSeatIDs(ctx context.Context, showtimeID uint, seatIDs []uint) ([]uint, error) {
	if len(seatIDs) == 0 {
		return []uint{}, nil
	}

	var bookedSeatIDs []uint

	err := r.DB.WithContext(ctx).Table("booking_seats AS bs").Select("bs.seat_id").Joins("JOIN bookings as b on b.id = bs.booking_id").Where("bs.showtime_id = ?", showtimeID).Where("bs.seat_id IN ?", seatIDs).Where("b.status = 'confirmed' OR (b.status = 'pending' AND b.expires_at > ?)", time.Now()).Pluck("bs.seat_id", &bookedSeatIDs).Error

	return bookedSeatIDs, err
}

func (r *BookingRepositoryImpl) CancelBooking(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Model(&models.Booking{}).Where("id = ? AND status = ?", id, "pending").Update("status", "cancelled")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrBookingNotFound
	}

	return nil
}

func (r *BookingRepositoryImpl) DeleteBookingSeats(ctx context.Context, tx *gorm.DB, bookingID uint) error {
	return tx.WithContext(ctx).Where("booking_id = ?", bookingID).Delete(&models.BookingSeat{}).Error
}

func (r *BookingRepositoryImpl) FindExpiredPendingBookings(ctx context.Context, now time.Time) ([]models.Booking, error) {
	var bookings []models.Booking

	err := r.DB.WithContext(ctx).Where("status = ?", "pending").Where("expires_at IS NOT NULL").Where("expires_at <= ?", now).Find(&bookings).Error

	return bookings, err
}

func (r *BookingRepositoryImpl) ExpireBooking(ctx context.Context, tx *gorm.DB, bookingID uint) error {
	result := tx.WithContext(ctx).Model(&models.Booking{}).Where("id = ?", bookingID).Where("status = ?", "pending").Update("status", "expired")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrBookingNotPending
	}

	return nil
}

func (r *BookingRepositoryImpl) FindBookingByIDForUpdate(ctx context.Context, id uint) (*models.Booking, error) {
	var booking models.Booking

	err := r.DB.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Preload("BookingSeats").Preload("BookingSeats.Seat").First(&booking, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepositoryImpl) ConfirmBooking(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Model(&models.Booking{}).Where("id = ?", id).Where("status = ?", "pending").Update("status", "confirmed")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrBookingNotPending
	}

	return nil
}
