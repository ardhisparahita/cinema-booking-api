package errors

import "errors"

var (
	ErrBookingNotFound     = errors.New("booking not found")
	ErrBookingCannotCancel = errors.New("booking cannot be cancelled")
	ErrBookingNotPending   = errors.New("booking is not pending")
	ErrBookingExpired      = errors.New("booking has expired")
)
