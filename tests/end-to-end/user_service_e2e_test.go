package endtoend

import (
	"bytes"
	"encoding/json"
	"flyhorizons-userservice/models"
	"flyhorizons-userservice/models/enums"
	"flyhorizons-userservice/repositories"
	entities "flyhorizons-userservice/repositories/entity"
	"flyhorizons-userservice/routes"
	"flyhorizons-userservice/services"
	"flyhorizons-userservice/services/authentication"
	"flyhorizons-userservice/services/converter"
	"flyhorizons-userservice/services/validation"
	mock_repositories "flyhorizons-userservice/tests/mocks"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type UserServiceEndToEndTests struct {
	repositories.BaseRepository
}

// Create a test version of BaseRepository that uses an in-memory SQLite database
func (repo *UserServiceEndToEndTests) CreateConnection() (*gorm.DB, error) {
	if repo.DB != nil {
		return repo.DB, nil
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to create SQLite database: %v", err)
		return nil, err
	}

	// Auto migrate entities for the test database
	err = db.AutoMigrate(&entities.UserEntity{})
	if err != nil {
		return nil, err
	}

	repo.DB = db
	return db, nil
}

func NewTestUserRepository() *repositories.UserRepository {
	baseRepo := &UserServiceEndToEndTests{}
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{}) // No shared cache
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	// Auto-migrate tables for the test database
	if err := db.AutoMigrate(&entities.UserEntity{}); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	baseRepo.DB = db
	return repositories.NewUserRepository(&baseRepo.BaseRepository)
}

// Adds users to the database on every run
func setupUsers(repo *repositories.UserRepository) {
	// Users
	testUsers := []entities.UserEntity{
		{ID: 1, FullName: "John Doe", Email: "john@doe.it", AccountType: 1, Password: "$2a$12$XbjoIVKp5miCCKU87B83S.Z5/OUMjS7OyQ5pW.UoieAyUeFW2G4q2", CreatedAt: time.Date(2025, time.March, 31, 10, 30, 0, 0, time.UTC)},
		{ID: 2, FullName: "Jane Doe", Email: "jane@doe.nl", AccountType: 0, Password: "$2a$12$7/NMoWfjAzIZhkK/6S4yy.Pnvo1YbF1lxIR2ehNgKfz0xzVyExzZO", CreatedAt: time.Date(2025, time.March, 31, 10, 30, 0, 0, time.UTC)},
	}

	// Add users to the test database
	for _, user := range testUsers {
		createdUser := repo.Create(user)
		log.Printf("Created user: %+v", createdUser)
	}
}

// Setup
func setupUserService(repo *repositories.UserRepository) *services.UserService {
	userConverter := converter.UserConverter{}
	passwordValidator := validation.PasswordValidator{}
	accountHashing := authentication.NewAccountHashing()
	return services.NewUserService(repo, accountHashing, passwordValidator, userConverter)
}

func setupUserRouter(service services.UserService, gatewayAuthMiddleware *mock_repositories.MockGatewayAuthMiddleware) *gin.Engine {
	router := gin.Default()
	routes.RegisterUserRoutes(router, &service, gatewayAuthMiddleware)
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

// End-to-End Tests
func TestEndToEndGetAllAsAdminReturnsUsers(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("admin", 1)
	bearerToken := "Bearer mocktoken12345"
	mockUsers := getUsers()
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)

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
}

func TestEndToEndGetAllAsUserReturnsAccessDenied(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	bearerToken := "Bearer mocktoken12345"
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)

	url := "/users/"
	httpRequest, _ := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestEndToEndGetUserByMatchingIDReturnsUser(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	bearerToken := "Bearer mocktoken12345"
	mockUser := getUsers()[0]
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)

	userID := 1
	url := fmt.Sprint("/users/", userID)
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
}

func TestEndToEndGetUserByNonMatchingIDReturnsUserNotFoundError(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	bearerToken := "Bearer mocktoken12345"
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)

	userID := 999
	url := fmt.Sprint("/users/", userID)
	httpRequest, _ := http.NewRequest("GET", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestEndToEndCreateNonExistingUserReturnsNewUser(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockUser := models.User{
		FullName:    "Clark Kent",
		Email:       "clark@torch.com",
		AccountType: enums.User,
		Password:    "HashedPasswordABC1234!",
	}
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)
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
	assert.Equal(t, mockUser.FullName, user.FullName)
	assert.Equal(t, mockUser.AccountType, user.AccountType)
	assert.Equal(t, mockUser.Email, user.Email)
}

func TestEndToEndCreateExistingUserReturnsUserExistsError(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	mockUser := getUsers()[0]
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)
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

func TestEndToEndDeleteUserByNonMatchingIDReturnsAccessDenied(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	bearerToken := "Bearer mocktoken12345"
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)

	userID := 2
	url := fmt.Sprintf("/users/%d", userID)
	httpRequest, _ := http.NewRequest("DELETE", url, nil)
	httpRequest.Header.Set("Authorization", bearerToken)

	responseRecorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(responseRecorder, httpRequest)

	// Assert
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}

func TestEndToEndUpdateUserByMatchingIDReturnsUpdatedUser(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 1)
	bearerToken := "Bearer mocktoken12345"
	mockUser := getUsers()[0]
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)
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
	assert.Equal(t, mockUser.ID, user.ID)
	assert.Equal(t, mockUser.FullName, user.FullName)
	assert.Equal(t, mockUser.AccountType, user.AccountType)
	assert.Equal(t, mockUser.Email, user.Email)
}

func TestEndToEndUpdateUserByNonMatchingIDReturnsAccessDenied(t *testing.T) {
	// Arrange
	// Setup repository
	userRepo := NewTestUserRepository()
	setupUsers(userRepo)
	// Setup service
	userService := setupUserService(userRepo)
	mockAPIGatewayMiddleware := mock_repositories.NewMockGatewayAuthMiddleware("user", 2)
	bearerToken := "Bearer mocktoken12345"
	mockUser := getUsers()[0]
	// Setup router
	router := setupUserRouter(*userService, mockAPIGatewayMiddleware)
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
}
