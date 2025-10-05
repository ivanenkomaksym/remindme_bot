package errors

import "fmt"

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Message, e.Err.Error(), e.Code)
	}
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Predefined domain errors
var (
	ErrUserNotFound = &DomainError{
		Code:    "USER_NOT_FOUND",
		Message: "User not found",
	}

	ErrUserExists = &DomainError{
		Code:    "USER_EXISTS",
		Message: "User already exists",
	}

	ErrUserDataInvalid = &DomainError{
		Code:    "INVALID_USER_DATA",
		Message: "User data invalid",
	}

	ErrReminderNotFound = &DomainError{
		Code:    "REMINDER_NOT_FOUND",
		Message: "Reminder not found",
	}

	ErrInvalidRecurrenceType = &DomainError{
		Code:    "INVALID_RECURRENCE_TYPE",
		Message: "Invalid recurrence type",
	}

	ErrInvalidTimeFormat = &DomainError{
		Code:    "INVALID_TIME_FORMAT",
		Message: "Invalid time format",
	}

	ErrEmptyMessage = &DomainError{
		Code:    "EMPTY_MESSAGE",
		Message: "Reminder message cannot be empty",
	}

	ErrUnauthorized = &DomainError{
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized access",
	}
)

// NewDomainError creates a new domain error
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
