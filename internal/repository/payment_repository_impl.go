package repository

import (
	"context"
	"errors"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
)

type PaymentRepositoryImpl struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{
		DB: db,
	}
}

func (r *PaymentRepositoryImpl) CreatePayment(ctx context.Context, tx *gorm.DB, payment *models.Payment) error {
	return tx.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepositoryImpl) FindPaymentByID(ctx context.Context, id uint) (*models.Payment, error) {
	var payment models.Payment

	err := r.DB.WithContext(ctx).Preload("Booking").Preload("Booking.BookingSeats").Preload("Booking.BookingSeats.Seats").First(&payment, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}
func (r *PaymentRepositoryImpl) FindPaymentByBookingID(ctx context.Context, bookingID uint) (*models.Payment, error) {
	var payment models.Payment

	err := r.DB.WithContext(ctx).Where("booking_id = ?", bookingID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepositoryImpl) UpdatePayment(ctx context.Context, tx *gorm.DB, payment *models.Payment) error {
	result := tx.WithContext(ctx).Model(&models.Payment{}).Where("id = ?", payment.ID).Where("status = ?", "pending").Updates(map[string]any{
		"status":       payment.Status,
		"provider_ref": payment.ProviderRef,
		"paid_at":      payment.PaidAt,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
