package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("resource already exists")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrUnauthorized      = errors.New("unauthorized")
)
