package validation_test

import (
	"flyhorizons-userservice/services/errors"
	"flyhorizons-userservice/services/validation"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPasswordValidator struct {
}

// Setup
func setup() validation.PasswordValidator {
	return validation.PasswordValidator{}
}

// Validator Tests
func TestValidateSecurePasswordReturnsNil(t *testing.T) {
	// Arrange
	passwordValidator := setup()
	password := "Fontysict1234!"

	// Act
	err := passwordValidator.Validate(password)

	// Assert
	assert.NoError(t, err)
}

func TestValidateInsufficientLengthPasswordThrowsException(t *testing.T) {
	// Arrange
	passwordValidator := setup()
	password := "FONTYS"

	// Act
	err := passwordValidator.Validate(password)

	// Assert
	assert.Error(t, err)
	_, ok := err.(*errors.InsufficientPasswordLengthError)
	assert.True(t, ok)
}

func TestValidateInconsistentPasswordThrowsException(t *testing.T) {
	// Arrange
	passwordValidator := setup()
	password := "lowercasepassword"

	// Act
	err := passwordValidator.Validate(password)

	// Assert
	assert.Error(t, err)
	_, ok := err.(*errors.InsufficientPasswordComplexityError)
	assert.True(t, ok)
}

func TestValidateCommonPasswordThrowsException(t *testing.T) {
	// Arrange
	passwordValidator := setup()
	password := "superman"

	// Act
	err := passwordValidator.Validate(password)

	// Assert
	assert.Error(t, err)
	_, ok := err.(*errors.CommonPasswordError)
	assert.False(t, ok)
}
