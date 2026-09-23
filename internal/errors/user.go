package errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrRefreshNotFound    = errors.New("refresh token not found")
	ErrRefreshExpired     = errors.New("refresh token expired")
	ErrRefreshRevoked     = errors.New("refresh token revoked")
)
