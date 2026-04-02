package model

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("requested item is not found")
	ErrConflict            = errors.New("item already exists")
	ErrBadParamInput       = errors.New("given param is not valid")
)
