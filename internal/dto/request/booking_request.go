package request

type CreateBookingRequest struct {
	ShowtimeID uint   `json:"showtime_id" validate:"required"`
	SeatIDs    []uint `json:"seat_ids" validate:"required,min=1"`
}
