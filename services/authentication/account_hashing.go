package authentication

import (
	"flyhorizons-userservice/services/errors"

	"golang.org/x/crypto/bcrypt"
)

type AccountHashing struct{}

func NewAccountHashing() *AccountHashing {
	return &AccountHashing{}
}

func (service *AccountHashing) HashPassword(password string) (string, error) {
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Return the custom hashing error
		return "", errors.NewHashingPasswordError(100)
	}
	return string(hashedPassword), nil
}
