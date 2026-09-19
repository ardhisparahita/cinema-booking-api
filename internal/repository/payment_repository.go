package repository

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, tx *gorm.DB, payment *models.Payment) error
	FindPaymentByID(ctx context.Context, id uint) (*models.Payment, error)
	FindPaymentByBookingID(ctx context.Context, bookingID uint) (*models.Payment, error)
	UpdatePayment(ctx context.Context, tx *gorm.DB, payment *models.Payment) error
}
