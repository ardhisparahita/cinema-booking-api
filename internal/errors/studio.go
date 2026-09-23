package errors

import "errors"

var (
	ErrStudioNotFound      = errors.New("studio not found")
	ErrStudioAlreadyExists = errors.New("studio already exists")
	ErrTheaterNotFound     = errors.New("theater not found")
)