package errors

import "fmt"

type InvalidCredentialsError struct {
	ErrorCode int
}

func (e *InvalidCredentialsError) Error() string {
	return fmt.Sprintf("The credentials provided are invalid. Error Code: %d", e.ErrorCode)
}

func NewInvalidCredentialsError(errorCode int) *InvalidCredentialsError {
	return &InvalidCredentialsError{ErrorCode: errorCode}
}
