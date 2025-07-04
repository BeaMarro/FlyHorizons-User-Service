package routes_test

import (
	"bytes"
	"encoding/json"
	"flyhorizons-userservice/models"
	"flyhorizons-userservice/routes"
	"flyhorizons-userservice/services/errors"
	mock_repositories "flyhorizons-userservice/tests/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestUserRoute struct {
}

// Setup
func setupUserRouter(mockService *mock_repositories.MockUserService, gatewayAuthMiddleware *mock_repositories.MockGatewayAuthMiddleware) *gin.Engine {
	router := gin.Default()

	routes.RegisterUserRoutes(router, mockService, gatewayAuthMiddleware)

	return router
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

// Router Integration Tests
func TestGetAllAsAdminReturnsUsersJSON(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	bearerToken := "Bearer mocktoken12345"
	mockUsers := getUsers()
	mockService.On("GetAll").Return(mockUsers)

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	url := "/users/"
	httpRequest, _ := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Unmarshal the JSON response
	var users []models.User
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &users)

	assert.NoError(t, err)
	assert.Equal(t, mockUsers, users)
	mockService.AssertExpectations(t)
}

func TestGetAllAsNonAdminRoleReturnsAccessDenied(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	bearerToken := "Bearer mocktoken12345"

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	url := "/users/"
	httpRequest, _ := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestGetByValidIDUsingLoggedInIDReturnsUserJSON(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	bearerToken := "Bearer mocktoken12345"
	mockUser := getUsers()[0]
	mockUserID := mockUser.ID
	userPtr := &mockUser
	mockService.On("GetByID", mockUserID).Return(userPtr, nil)

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/users/%d", mockUserID)
	httpRequest, _ := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// Unmarshal the JSON response
	var user models.User
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &user)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
	mockService.AssertExpectations(t)
}

func TestGetByInvalidIDUsingLoggedInIDReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 999)
	bearerToken := "Bearer mocktoken12345"
	mockUser := getUsers()[0]
	mockUserID := mockUser.ID
	userPtr := &mockUser
	mockService.On("GetByID", mockUserID).Return(userPtr, nil)

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/users/%d", mockUserID)
	httpRequest, _ := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestCreateNonExistingUserReturnsCreatedUser(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := new(mock_repositories.MockGatewayAuthMiddleware)
	mockUser := getUsers()[0]
	mockService.On("Create", mockUser).Return(&mockUser, nil)

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the user
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                        // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	var user models.User
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
	mockService.AssertExpectations(t)
}

func TestCreateExistingUserReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := new(mock_repositories.MockGatewayAuthMiddleware)
	mockUser := getUsers()[0]
	errorCode := 409
	mockService.On("Create", mockUser).Return(nil, errors.NewUserExistsError(mockUser.ID, errorCode))

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the user
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                        // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
}

func TestCreateUserWithUnsufficientPasswordLengthReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := new(mock_repositories.MockGatewayAuthMiddleware)
	mockUser := getUsers()[0]
	errorCode := 400
	mockService.On("Create", mockUser).Return(nil, errors.NewInsufficientPasswordLengthError(errorCode))

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the user
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                        // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestCreateUserWithInsufficientPasswordComplexityReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := new(mock_repositories.MockGatewayAuthMiddleware)
	mockUser := getUsers()[0]
	errorCode := 400
	mockService.On("Create", mockUser).Return(nil, errors.NewInsufficientPasswordComplexityError(errorCode))

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the user
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                        // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestCreateUserWithCommonPasswordReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := new(mock_repositories.MockGatewayAuthMiddleware)
	mockUser := getUsers()[0]
	errorCode := 400
	mockService.On("Create", mockUser).Return(nil, errors.NewCommonPasswordError(errorCode))

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the user
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody)) // JSON body
	httpRequest.Header.Set("Content-Type", "application/json")                        // Set the Content-Type header
	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestDeleteExistingUserAsMatchingRoleReturnsHTTPStatusOk(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	userID := 1
	bearerToken := "Bearer mocktoken12345"
	mockService.On("DeleteByID", userID).Return(true, nil)

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/users/%d", userID)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteNonExistingUserAsMatchingRoleReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 999)
	bearerToken := "Bearer mocktoken12345"
	userID := 1
	mockService.On("DeleteByID", userID).Return(false, errors.NewUserNotFoundError(userID, 404))

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	url := fmt.Sprintf("/users/%d", userID)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
}

func TestUpdateExistingUserAsMatchingRoleReturnsUpdatedUser(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockUser := getUsers()[0]
	userPtr := &mockUser
	mockService.On("Update", mockUser).Return(userPtr, nil)
	bearerToken := "Bearer mocktoken12345"

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the airport
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("PUT", "/users/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var user models.User
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
	mockService.AssertExpectations(t)
}

func TestUpdateUserWithUnsufficientPasswordLengthReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockUser := getUsers()[0]
	mockService.On("Update", mockUser).Return(nil, errors.NewInsufficientPasswordLengthError(400))
	bearerToken := "Bearer mocktoken12345"

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the airport
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("PUT", "/users/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestUpdateUserWithUnsufficientPasswordComplexityReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockUser := getUsers()[0]
	mockService.On("Update", mockUser).Return(nil, errors.NewInsufficientPasswordComplexityError(400))
	bearerToken := "Bearer mocktoken12345"

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the airport
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("PUT", "/users/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestUpdateUserWithCommonPasswordReturnsHTTPStatusError(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	mockUser := getUsers()[0]
	mockService.On("Update", mockUser).Return(nil, errors.NewCommonPasswordError(400))
	bearerToken := "Bearer mocktoken12345"

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the airport
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("PUT", "/users/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestUpdateNonMatchingUserReturnsAccessDenied(t *testing.T) {
	// Arrange
	mockService := new(mock_repositories.MockUserService)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 999)
	mockUser := getUsers()[0]
	userPtr := &mockUser
	mockService.On("Update", mockUser).Return(userPtr, nil)
	bearerToken := "Bearer mocktoken12345"

	router := setupUserRouter(mockService, mockAPIGatewayMiddleware)

	// Make the JSON to create the airport
	requestBody, _ := json.Marshal(mockUser)
	httpRequest, _ := http.NewRequest("PUT", "/users/", bytes.NewBuffer(requestBody))
	httpRequest.Header.Set("Content-Type", "application/json") // Set the Content-Type header
	httpRequest.Header.Set("Authorization", bearerToken)       // Set the Bearer token

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)

	var errResponse map[string]interface{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
}
