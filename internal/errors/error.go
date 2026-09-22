package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenInput  = errors.New("invalid token input")
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

var (
	ErrBookingNotFound     = errors.New("booking not found")
	ErrBookingCannotCancel = errors.New("booking cannot be cancelled")
	ErrBookingNotPending   = errors.New("booking is not pending")
	ErrBookingExpired      = errors.New("booking has expired")
)

var (
	ErrShowtimeNotFound      = errors.New("showtime not found")
	ErrShowtimeFinished      = errors.New("showtime already finished")
	ErrInvalidShowtimeTime   = errors.New("end time after must be after start time")
	ErrShowtimeConflict      = errors.New("showtime conflict with another showtime")
	ErrInvalidShowtimePrice  = errors.New("price must be greater than zero")
	ErrInvalidShowtimeInPast = errors.New("showtime cannot be in the past")
)
