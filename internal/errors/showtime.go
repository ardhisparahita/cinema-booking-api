package errors

import "errors"

var (
	ErrShowtimeNotFound      = errors.New("showtime not found")
	ErrShowtimeFinished      = errors.New("showtime already finished")
	ErrInvalidShowtimeTime   = errors.New("end time after must be after start time")
	ErrShowtimeConflict      = errors.New("showtime conflict with another showtime")
	ErrInvalidShowtimePrice  = errors.New("price must be greater than zero")
	ErrInvalidShowtimeInPast = errors.New("showtime cannot be in the past")
)
