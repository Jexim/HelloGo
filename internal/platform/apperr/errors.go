package apperr

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
	ErrBadRequest    = errors.New("bad request")
)
