package errors

import "fmt"

type CommonPasswordError struct {
	ErrorCode int
}

func (e *CommonPasswordError) Error() string {
	return fmt.Sprintf("The password belongs to the list of most common passwords, thus it is not sufficiently secure. [Error code: %d]", e.ErrorCode)
}

func NewCommonPasswordError(errorCode int) *CommonPasswordError {
	return &CommonPasswordError{ErrorCode: errorCode}
}
