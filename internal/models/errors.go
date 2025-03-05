package models

import (
	"errors"
	"fmt"
)
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

type AuthError struct {
	Code    string
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("auth error [%s]: %s", e.Code, e.Message)
}

type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

var (
	ErrUserNotFound        = &AuthError{Code: "USER_001", Message: "User not found"}
	ErrInvalidCredentials  = &AuthError{Code: "AUTH_001", Message: "Invalid email or password"}
	ErrUserAlreadyExists   = &AuthError{Code: "USER_002", Message: "User with this email already exists"}
	ErrUnauthorized        = &AuthError{Code: "AUTH_002", Message: "Unauthorized access"}

	ErrBoardNotFound       = errors.New("board not found")
	ErrInsufficientAccess  = errors.New("insufficient access rights")

	ErrColumnNotFound      = errors.New("column not found")

	ErrCardNotFound        = errors.New("card not found")

	ErrCommentNotFound     = errors.New("comment not found")

	ErrLabelNotFound       = errors.New("label not found")
)

func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

func IsAuthError(err error) bool {
	_, ok := err.(*AuthError)
	return ok
}

func IsDatabaseError(err error) bool {
	_, ok := err.(*DatabaseError)
	return ok
}

func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func NewAuthError(code, message string) error {
	return &AuthError{
		Code:    code,
		Message: message,
	}
}

func NewDatabaseError(operation string, err error) error {
	return &DatabaseError{
		Operation: operation,
		Err:       err,
	}
}