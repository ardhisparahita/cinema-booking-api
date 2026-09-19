package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, userID uint, req *request.CreatePaymentRequest) (*response.PaymentResponse, error)
	GetPaymentByID(ctx context.Context, userID uint, id uint) (*response.PaymentResponse, error)
	ConfirmPayment(ctx context.Context, userID uint, id uint) (*response.PaymentResponse, error)
}
