package errors

import "errors"

var (
	ErrPostNotFound  = errors.New("post not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrDuplicatePost = errors.New("duplicate post")
	ErrInternal      = errors.New("internal error")
)
