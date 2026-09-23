package errors

import "errors"

var (
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrPaymentAlreadyExist      = errors.New("payment already exist")
	ErrPaymentAlreadyPaid       = errors.New("payment already paid")
	ErrInvalidPaymentMethod     = errors.New("invalid payment method")
	ErrPaymentCannotConfirm     = errors.New("payment cannot be confirmed")
)
