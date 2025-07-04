package services_test

import (
	"flyhorizons-userservice/models"
	entities "flyhorizons-userservice/repositories/entity"
	"flyhorizons-userservice/services"
	"flyhorizons-userservice/services/authentication"
	"flyhorizons-userservice/services/converter"
	"flyhorizons-userservice/services/errors"
	"flyhorizons-userservice/services/validation"
	mock_repositories "flyhorizons-userservice/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestUserService struct {
}

// Setup
func setupUserService() (*mock_repositories.MockUserRepository, *services.UserService) {
	mockRepo := new(mock_repositories.MockUserRepository)
	accountHashing := new(authentication.AccountHashing)
	userConverter := new(converter.UserConverter)
	passwordValidator := new(validation.PasswordValidator)
	userService := services.NewUserService(mockRepo, accountHashing, *passwordValidator, *userConverter)
	return mockRepo, userService
}

func getCurrentDateTime() time.Time {
	return time.Now()
}

func getUserEntities() []entities.UserEntity {
	return []entities.UserEntity{
		{
			ID:          1,
			FullName:    "John Doe",
			Email:       "john@doe.it",
			AccountType: 1,
			Password:    "$2a$12$XbjoIVKp5miCCKU87B83S.Z5/OUMjS7OyQ5pW.UoieAyUeFW2G4q2", // 1234!
			CreatedAt:   getCurrentDateTime(),
		},
		{
			ID:          2,
			FullName:    "Jane Doe",
			Email:       "jane@doe.nl",
			AccountType: 0,
			Password:    "$2a$12$7/NMoWfjAzIZhkK/6S4yy.Pnvo1YbF1lxIR2ehNgKfz0xzVyExzZO", // 4321!
			CreatedAt:   getCurrentDateTime(),
		},
	}
}

func getUsers() []models.User {
	return []models.User{
		{
			ID:          1,
			FullName:    "John Doe",
			Email:       "john@doe.it",
			AccountType: 1,
			Password:    "$2a$12$XbjoIVKp5miCCKU87B83S.Z5/OUMjS7OyQ5pW.UoieAyUeFW2G4q2", // 1234!
		},
		{
			ID:          2,
			FullName:    "Jane Doe",
			Email:       "jane@doe.nl",
			AccountType: 0,
			Password:    "$2a$12$7/NMoWfjAzIZhkK/6S4yy.Pnvo1YbF1lxIR2ehNgKfz0xzVyExzZO", // 4321!
		},
	}
}

// Service Unit Tests
func TestGetAllUsersReturnsUsers(t *testing.T) {
	// Arrange
	mockRepo, userService := setupUserService()
	userEntities := getUserEntities()
	expected := getUsers()
	mockRepo.On("GetAll").Return(userEntities)

	// Act
	users := userService.GetAll()

	// Assert
	assert.Equal(t, expected, users)
}

func TestGetUserByValidIDReturnsUser(t *testing.T) {
	// Arrange
	mockRepo, userService := setupUserService()
	userEntity := getUserEntities()[0]
	userID := 1
	mockRepo.On("GetByID", userID).Return(userEntity)
	expected := getUsers()[0]

	// Act
	user, err := userService.GetByID(userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, *user)
}

func TestGetByInvalidIDThrowsException(t *testing.T) {
	// Arrange
	mockRepo, userService := setupUserService()
	userID := 999
	mockRepo.On("GetByID", userID).Return(entities.UserEntity{}) // Returns an empty UserEntity

	// Act
	user, err := userService.GetByID(userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.NewUserNotFoundError(userID, 404), err)
	assert.Nil(t, user)
}

func TestCreateNonExistingUserReturnsCreatedUser(t *testing.T) {
	// Arrange
	mockRepo, userService := setupUserService()
	user := getUsers()[0]
	userEntity := getUserEntities()[0]

	mockRepo.On("GetAll").Return([]entities.UserEntity{getUserEntities()[1]})
	mockRepo.On("Create", mock.MatchedBy(func(u entities.UserEntity) bool {
		return u.ID == userEntity.ID // Ignore password and CreatedAt differences
	})).Return(userEntity)

	// Act
	postUser, err := userService.Create(user)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.ID, postUser.ID)
	assert.Equal(t, user.FullName, postUser.FullName)
	assert.Equal(t, user.Email, postUser.Email)
	assert.Equal(t, user.AccountType, postUser.AccountType)
}

func TestCreateExistingUserThrowsException(t *testing.T) {
	// Arrange
	mockRepo, userService := setupUserService()
	user := getUsers()[0]
	userEntity := getUserEntities()[0]

	mockRepo.On("GetAll").Return([]entities.UserEntity{getUserEntities()[0]})
	mockRepo.On("Create", mock.MatchedBy(func(u entities.UserEntity) bool {
		return u.ID == userEntity.ID // Ignore password and CreatedAt differences
	})).Return(userEntity)

	// Act
	postUser, err := userService.Create(user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.NewUserExistsError(1, 409), err)
	assert.Nil(t, postUser)
}

// TODO: Fix
// This fails after the RabbitMQ implementation, probably the RabbitMQ shall be mocked here
// func TestDeleteByExistingIDReturnsTrue(t *testing.T) {
// 	// Arrange
// 	mockRepo, userService := setupUserService()
// 	userID := 1
// 	mockRepo.On("GetAll").Return([]entities.UserEntity{getUserEntities()[0]})
// 	mockRepo.On("DeleteByID", userID).Return(true)

// 	// Act
// 	isDeleted, err := userService.DeleteByID(userID)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.True(t, isDeleted)
// }

func TestDeleteByNonExistingIDThrowsException(t *testing.T) {
	// Arrange
	mockRepo, userService := setupUserService()
	userID := 999
	mockRepo.On("GetAll").Return([]entities.UserEntity{getUserEntities()[1]})
	mockRepo.On("DeleteByID", userID).Return(false)

	// Act
	isDeleted, err := userService.DeleteByID(userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.NewUserNotFoundError(userID, 404), err)
	assert.False(t, isDeleted)
}

func TestUpdateByExistingUserReturnsUpdatedUser(t *testing.T) {
	// Arrange
	mockRepo, userService := setupUserService()
	user := getUsers()[0]
	userEntity := getUserEntities()[0]

	mockRepo.On("GetAll").Return(getUserEntities())
	mockRepo.On("Update", mock.MatchedBy(func(u entities.UserEntity) bool {
		return u.ID == userEntity.ID
	})).Return(userEntity)

	// Act
	putUser, err := userService.Update(user)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user.ID, putUser.ID)
	assert.Equal(t, user.FullName, putUser.FullName)
	assert.Equal(t, user.Email, putUser.Email)
	assert.Equal(t, user.AccountType, putUser.AccountType)
}

func TestUpdateByNonExistingUserReturnsUpdatedUser(t *testing.T) {
	// Arrange
	mockRepo, userService := setupUserService()
	user := getUsers()[0]
	userEntity := getUserEntities()[0]

	mockRepo.On("GetAll").Return([]entities.UserEntity{})
	mockRepo.On("Update", mock.MatchedBy(func(u entities.UserEntity) bool {
		return u.ID == userEntity.ID
	})).Return(userEntity)

	// Act
	putUser, err := userService.Update(user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.NewUserNotFoundError(user.ID, 404), err)
	assert.Nil(t, putUser)
}
