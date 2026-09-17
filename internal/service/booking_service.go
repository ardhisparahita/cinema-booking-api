package service

import (
	"context"

	"github.com/ardhisparahita/cinema-booking-api/internal/dto/request"
	"github.com/ardhisparahita/cinema-booking-api/internal/dto/response"
)

type BookingService interface {
	CreateBooking(ctx context.Context, userID uint, req request.CreateBookingRequest) (*response.BookingResponse, error)
	GetBookingByID(ctx context.Context, userID uint, id uint) (*response.BookingResponse, error)
	GetMyBookings(ctx context.Context, userID uint) ([]response.BookingResponse, error)
	CancelBooking(ctx context.Context, userID uint, id uint) error
}
