package domain

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	// ErrInvalidInput is returned when input data fails validation
	// Examples:
	//   - pack sizes are not unique
	//   - pack sizes <= 0 or > 1e6
	//   - amount <= 0
	ErrInvalidInput = errors.New("invalid input")

	// ErrNoSolutionStrict is returned when no exact solution without overage is found
	// This can happen in strict search mode when a solution
	// with zero overage is required (overage = 0)
	ErrNoSolutionStrict = errors.New("no strict solution found")

	// ErrNoSolution is returned when no solution is found at all
	// Usually this should not happen, as the smallest pack size
	// can always be used
	ErrNoSolution = errors.New("no solution found")

	// ErrPackSizeSetNotFound is returned when pack size set is not found
	ErrPackSizeSetNotFound = errors.New("pack size set not found")

	// ErrPackSizeSetAlreadyExists is returned when a set with this name already exists
	ErrPackSizeSetAlreadyExists = errors.New("pack size set already exists")

	// ErrSolutionNotFound is returned when solution is not found in cache
	ErrSolutionNotFound = errors.New("solution not found in cache")

	// ErrCacheUnavailable is returned when cache is unavailable
	ErrCacheUnavailable = errors.New("cache unavailable")
)

// ValidationError represents a validation error with additional context
type ValidationError struct {
	Field   string // Field that failed validation
	Value   any    // Value that failed validation
	Message string // Error message
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

// Unwrap allows using errors.Is and errors.As
func (e *ValidationError) Unwrap() error {
	return ErrInvalidInput
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value any, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// SolverError represents an error when solving the packing problem
type SolverError struct {
	Sizes   []int  // Pack sizes
	Amount  int    // Required amount
	Message string // Error message
	Err     error  // Internal error
}

// Error implements the error interface
func (e *SolverError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("solver error: %s (sizes: %v, amount: %d): %v", e.Message, e.Sizes, e.Amount, e.Err)
	}
	return fmt.Sprintf("solver error: %s (sizes: %v, amount: %d)", e.Message, e.Sizes, e.Amount)
}

// Unwrap allows using errors.Is and errors.As
func (e *SolverError) Unwrap() error {
	return e.Err
}

// NewSolverError creates a new solver error
func NewSolverError(sizes []int, amount int, message string, err error) *SolverError {
	return &SolverError{
		Sizes:   sizes,
		Amount:  amount,
		Message: message,
		Err:     err,
	}
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsNotFoundError checks if the error is a "not found" error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrPackSizeSetNotFound) || errors.Is(err, ErrSolutionNotFound)
}

// IsNoSolutionError checks if the error is a "solution not found" error
func IsNoSolutionError(err error) bool {
	return errors.Is(err, ErrNoSolution) || errors.Is(err, ErrNoSolutionStrict)
}
