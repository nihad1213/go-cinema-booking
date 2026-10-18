package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrUnauthorized    = errors.New("missing X-User-ID header")
	ErrSeatUnavailable = errors.New("one or more seats are not available")
)
