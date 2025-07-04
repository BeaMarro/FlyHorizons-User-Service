package converter_test

import (
	"flyhorizons-userservice/models"
	entities "flyhorizons-userservice/repositories/entity"
	"flyhorizons-userservice/services/converter"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestUserConverter struct {
}

// Setup
func setup() converter.UserConverter {
	return converter.UserConverter{}
}

func getCurrentDateTime() time.Time {
	return time.Now()
}

func getUserEntity() entities.UserEntity {
	return entities.UserEntity{
		ID:          2,
		FullName:    "Jane Doe",
		Email:       "jane@doe.nl",
		AccountType: 0,
		Password:    "$2a$12$7/NMoWfjAzIZhkK/6S4yy.Pnvo1YbF1lxIR2ehNgKfz0xzVyExzZO", // 4321!
		CreatedAt:   getCurrentDateTime(),
	}
}

func getUser() models.User {
	return models.User{
		ID:          2,
		FullName:    "Jane Doe",
		Email:       "jane@doe.nl",
		AccountType: 0,
		Password:    "$2a$12$7/NMoWfjAzIZhkK/6S4yy.Pnvo1YbF1lxIR2ehNgKfz0xzVyExzZO", // 4321!
	}
}

// Converter Tests
func TestConvertUserToUserEntityReturnsUserEntity(t *testing.T) {
	// Arrange
	userConverter := setup()
	user := getUser()

	// Act
	userEntity := userConverter.ConvertUserToUserEntity(user)

	// Assert
	assert.Equal(t, userEntity, getUserEntity())
}

func TestConvertUserEntityToUserReturnsUser(t *testing.T) {
	// Arrange
	userConverter := setup()
	userEntity := getUserEntity()

	// Act
	user := userConverter.ConvertUserEntityToUser(userEntity)

	// Assert
	assert.Equal(t, user, getUser())
}
