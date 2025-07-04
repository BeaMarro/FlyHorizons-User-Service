package errors

import "fmt"

type InvalidAccountTypeError struct {
	ErrorCode int
}

func (e *InvalidAccountTypeError) Error() string {
	return fmt.Sprintf("The account type provided is invalid. Error Code: %d", e.ErrorCode)
}

func NewInvalidAccountTypeError(errorCode int) *InvalidAccountTypeError {
	return &InvalidAccountTypeError{ErrorCode: errorCode}
}
