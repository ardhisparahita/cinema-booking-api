package repository

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, tx *gorm.DB, booking *models.Booking) error
	CreateBookingSeats(ctx context.Context, tx *gorm.DB, bookingSeats []models.BookingSeat) error
	FindBookingByID(ctx context.Context, id uint) (*models.Booking, error)
	FindBookingByUserID(ctx context.Context, userID uint) ([]models.Booking, error)
	FindBookedSeatIDs(ctx context.Context, showtimeID uint, seatIDs []uint) ([]uint, error)
	CancelBooking(ctx context.Context, id uint) error
	DeleteBookingSeats(ctx context.Context, tx *gorm.DB, bookingID uint) error
}
