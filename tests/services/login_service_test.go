package services_test

import (
	"flyhorizons-userservice/models/request"
	"flyhorizons-userservice/models/response"
	"flyhorizons-userservice/services"
	"flyhorizons-userservice/services/converter"
	"flyhorizons-userservice/services/errors"
	mock_repositories "flyhorizons-userservice/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestLoginService struct {
}

// Setup
func setupLoginService() (*mock_repositories.MockUserRepository, *mock_repositories.MockJwtTokenSigner, *services.LoginService) {
	mockRepo := new(mock_repositories.MockUserRepository)
	userConverter := new(converter.UserConverter)
	mockJwtTokenSigner := new(mock_repositories.MockJwtTokenSigner)
	loginService := services.NewLoginService(mockRepo, *userConverter, mockJwtTokenSigner)
	return mockRepo, mockJwtTokenSigner, loginService
}

func getLoginRequest(email string, password string) request.LoginRequest {
	return request.LoginRequest{
		Email:    email,
		Password: password,
	}
}

// Service Unit Tests
func TestLoginUsingCorrectCredentialsReturnsAccessToken(t *testing.T) {
	// Arrange
	mockRepo, mockJwtTokenSigner, loginService := setupLoginService()
	// User credentials
	id := 1
	email := "john@doe.it"
	password := "1234!"
	ip := "1234.123.12"
	mockAccessToken := "Mock Access Token"
	mockRepo.On("GetByEmail", email).Return(getUserEntities()[0])
	mockRepo.On("SaveLastLoginTime", id).Return()
	// Mock signing the Jwt auth token
	mockJwtTokenSigner.On("SignToken", mock.Anything).Return(mockAccessToken, nil)
	loginRequest := getLoginRequest(email, password)

	// Act
	accessToken, err := loginService.Login(loginRequest, ip)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, accessToken, &response.LoginResponse{
		AccessToken: mockAccessToken,
	})
}

func TestLoginUsingIncorrectCredentialsThrowsException(t *testing.T) {
	// Arrange
	mockRepo, mockJwtTokenSigner, loginService := setupLoginService()
	// Incorrect user credentials
	email := "john@doe.it"
	password := "4321!"
	ip := "1234.123.12"
	mockRepo.On("GetByEmail", email).Return(getUserEntities()[0])
	mockJwtTokenSigner.On("SignToken", mock.Anything).Return("", nil)
	loginRequest := getLoginRequest(email, password)

	// Act
	accessToken, err := loginService.Login(loginRequest, ip)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errors.NewInvalidCredentialsError(400), err)
	assert.Nil(t, accessToken)
}
