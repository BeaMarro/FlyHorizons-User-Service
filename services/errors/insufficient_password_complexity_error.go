package errors

import "fmt"

type InsufficientPasswordComplexityError struct {
	ErrorCode int
}

func (e *InsufficientPasswordComplexityError) Error() string {
	return fmt.Sprintf("The complexity given password is insufficient, it should consist of upper and lowercase letters, as well as a number and symbol. [Error code: %d]", e.ErrorCode)
}

func NewInsufficientPasswordComplexityError(errorCode int) *InsufficientPasswordComplexityError {
	return &InsufficientPasswordComplexityError{ErrorCode: errorCode}
}
