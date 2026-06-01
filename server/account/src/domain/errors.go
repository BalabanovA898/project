package domain

import "errors"
import "fmt"

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("resource already exists")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrUnauthorized      = errors.New("unauthorized")
)

// ValidationError представляет ошибку валидации поля
type ValidationError struct {
	Field   string
	Message string
}

// Error реализует интерфейс error
func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' - %s", ve.Field, ve.Message)
}

// NewValidationError создает новую ошибку валидации
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
