package errors

import "fmt"

type InsufficientPasswordLengthError struct {
	ErrorCode int
}

func (e *InsufficientPasswordLengthError) Error() string {
	return fmt.Sprintf("The length of the given password is insufficient, it should be at least 13 characters in length. [Error code: %d]", e.ErrorCode)
}

func NewInsufficientPasswordLengthError(errorCode int) *InsufficientPasswordLengthError {
	return &InsufficientPasswordLengthError{ErrorCode: errorCode}
}
