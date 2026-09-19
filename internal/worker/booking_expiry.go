package worker

import (
	"context"
	"log"
	"time"

	"github.com/ardhisparahita/cinema-booking-api/internal/service"
)

type BookingExpiryWorker struct {
	Service  service.BookingService
	Interval time.Duration
}

func NewBookingExpiryWorker(service service.BookingService, interval time.Duration) *BookingExpiryWorker {
	return &BookingExpiryWorker{
		Service:  service,
		Interval: interval,
	}
}

func (w *BookingExpiryWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.Service.ExpireBooking(ctx); err != nil {
				log.Printf("failed to expire bookings: %v", err)
			}
		case <-ctx.Done():
			return
		}

	}
}
