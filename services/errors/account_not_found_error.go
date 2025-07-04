package errors

import "fmt"

type UserNotFoundError struct {
	ID int
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("User with the ID %d was not found", e.ID)
}

func NewUserNotFoundError(id int, errorCode int) *UserNotFoundError {
	return &UserNotFoundError{ID: id}
}
