package apperrors

import "errors"

var (
	ErrInternalServer          = errors.New("internal server error")
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidLoginCredantials = errors.New("invalid email or password")
)
