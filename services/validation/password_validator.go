package validation

import (
	"bufio"
	"flyhorizons-userservice/services/errors"
	"fmt"
	"os"
	"unicode"
)

type PasswordValidator struct {
	commonPasswords map[string]struct{}
}

func NewPasswordValidator(filepath string) (*PasswordValidator, error) {
	validator := &PasswordValidator{
		commonPasswords: make(map[string]struct{}),
	}

	// Load common passwords list during initialization
	err := validator.LoadCommonPasswords(filepath)
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func (passwordValidator *PasswordValidator) LoadCommonPasswords(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open common passwords file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		password := scanner.Text()
		passwordValidator.commonPasswords[password] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading common passwords: %w", err)
	}

	return nil
}

// Internal methods
func (passwordValidator *PasswordValidator) CheckLength(password string) bool {
	return len(password) > 12
}

func (passwordValidator *PasswordValidator) CheckCharacterComplexity(password string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (passwordValidator *PasswordValidator) CheckAgainstWeakestPasswordsList(password string) bool {
	_, found := passwordValidator.commonPasswords[password]
	return found
}

// Validation method
func (passwordValidator *PasswordValidator) Validate(password string) error {
	if passwordValidator.CheckAgainstWeakestPasswordsList(password) {
		return errors.NewCommonPasswordError(400)
	}
	if !passwordValidator.CheckLength(password) {
		return errors.NewInsufficientPasswordLengthError(400)
	}
	if !passwordValidator.CheckCharacterComplexity(password) {
		return errors.NewInsufficientPasswordComplexityError(400)
	}
	return nil
}
