package errors

import "errors"

var (
	ErrSeatNotFound            = errors.New("seat not found")
	ErrSeatsNotFound           = errors.New("one or more seats not found")
	ErrSeatAlreadyExists       = errors.New("seat already exists")
	ErrSeatAlreadyBooked       = errors.New("one or more seats already booked")
	ErrSeatWrongStudio         = errors.New("one or more seats do not belong to showtime studio")
	ErrDuplicateSeat           = errors.New("duplicate seat in booking")
	ErrInvalidBookingSeats     = errors.New("at least one seat is required")
	ErrInvalidSeatType         = errors.New("invalid seat type")
	ErrInvalidSeatRow          = errors.New("invalid seat row")
	ErrInvalidSeatColumn       = errors.New("invalid seat column")
	ErrSeatOutsideStudioLayout = errors.New("seat position is outside studio layout")
	ErrSeatLocked              = errors.New("one or more seats are currently locked")
)
