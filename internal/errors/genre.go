package errors

import "errors"

var (
	ErrGenreNotFound      = errors.New("one or more genre not found")
	ErrGenreAlreadyExists = errors.New("genre already exists")
)
