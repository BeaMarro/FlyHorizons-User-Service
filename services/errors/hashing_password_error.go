package errors

import "fmt"

type HashingPasswordError struct {
	ErrorCode int
}

func (e *HashingPasswordError) Error() string {
	return fmt.Sprintf("Password Hashing Error - Code: %d", e.ErrorCode)
}

func NewHashingPasswordError(errorCode int) *HashingPasswordError {
	return &HashingPasswordError{ErrorCode: errorCode}
}
