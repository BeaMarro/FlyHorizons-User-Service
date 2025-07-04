package errors

import "fmt"

type UserExistsError struct {
	ID int
}

func (e *UserExistsError) Error() string {
	return fmt.Sprintf("User with the ID %d already exists", e.ID)
}

func NewUserExistsError(id int, errorCode int) *UserExistsError {
	return &UserExistsError{ID: id}
}
