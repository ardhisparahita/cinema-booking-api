package errors

import "errors"

var (
	ErrMovieNotFound      = errors.New("movie not found")
	ErrInvalidMovieRating = errors.New("invalid movie rating")
)