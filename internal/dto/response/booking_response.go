package response

import "time"

type BookingSeatResponse struct {
	ID         uint    `json:"id"`
	SeatID     uint    `json:"seat_id"`
	ShowtimeID uint    `json:"showtime_id"`
	Price      float64 `json:"price"`
}

type BookingResponse struct {
	ID           uint                  `json:"id"`
	BookingCode  string                `json:"booking_code"`
	UserID       uint                  `json:"user_id"`
	ShowtimeID   uint                  `json:"showtime_id"`
	TotalPrice   float64               `json:"total_price"`
	Status       string                `json:"status"`
	ExpiresAt    time.Time             `json:"expires_at"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	BookingSeats []BookingSeatResponse `json:"booking_seats"`
}
