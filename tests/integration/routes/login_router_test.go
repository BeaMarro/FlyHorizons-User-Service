package routes_test

import (
	"bytes"
	"encoding/json"
	"flyhorizons-userservice/models/request"
	"flyhorizons-userservice/models/response"
	"flyhorizons-userservice/routes"
	"flyhorizons-userservice/services/errors"
	mock_repositories "flyhorizons-userservice/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type AuthRouterTest struct {
}

// Setup
func setupLoginRouter(mockService *mock_repositories.MockLoginService, gatewayAuthMiddleware *mock_repositories.MockGatewayAuthMiddleware) *gin.Engine {
	router := gin.Default()

	routes.RegisterAuthRoutes(router, mockService)

	return router
}

func getCorrectLoginCredentials() request.LoginRequest {
	return request.LoginRequest{
		Email:    "john@doe.it",
		Password: "1234",
	}
}

func getIncorrectLoginCredentials() request.LoginRequest {
	return request.LoginRequest{
		Email:    "johathan@doe.it",
		Password: "12344",
	}
}

// Router Integration Tests
func TestLoginUsingCorrectCredentialsReturnsAccessToken(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockLoginService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockLoginRequest := getCorrectLoginCredentials()
	mockAccessToken := "Access-Token-Mock-1234"
	mockLoginResponse := &response.LoginResponse{AccessToken: mockAccessToken}
	mockService.On("Login", mockLoginRequest).Return(mockLoginResponse, nil)

	router := setupLoginRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the login request
	requestBody, _ := json.Marshal(mockLoginRequest)
	httpRequest, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	var responseBody response.LoginResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, mockAccessToken, responseBody.AccessToken)
	mockService.AssertExpectations(t)
}

func TestLoginUsingIncorrectCredentialsReturnsAccessDenied(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockLoginService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockLoginRequest := getIncorrectLoginCredentials()
	mockService.On("Login", mockLoginRequest).Return(nil, errors.NewInvalidCredentialsError(400))

	router := setupLoginRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the login request
	requestBody, _ := json.Marshal(mockLoginRequest)
	httpRequest, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)

	mockService.AssertExpectations(t)
}
