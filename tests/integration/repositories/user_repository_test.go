package repositories_test

import (
	"flyhorizons-userservice/repositories"
	entities "flyhorizons-userservice/repositories/entity"
	"log"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Create a test version of BaseRepository that uses an in-memory SQLite database
type TestBaseRepository struct {
	repositories.BaseRepository
}

func (repo *TestBaseRepository) CreateConnection() (*gorm.DB, error) {
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
	baseRepo := &TestBaseRepository{}
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
func setupUsers(repo *repositories.UserRepository) []entities.UserEntity {
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

	return testUsers
}

// Integration Database Tests
func TestUserRepositoryGetAllReturnsUsers(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	testUsers := setupUsers(userRepo)

	// Act
	users := userRepo.GetAll()

	// Assert
	assert.Equal(t, testUsers, users)
}

func TestUserRepositoryGetByValidIDReturnsUser(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	testUsers := setupUsers(userRepo)
	userID := 1

	// Act
	user := userRepo.GetByID(userID)

	// Assert
	assert.Equal(t, testUsers[0].ID, user.ID)
}

func TestUserRepositoryGetByInvalidIDReturnsEmptyUser(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	invalidUserID := 999

	// Act
	user := userRepo.GetByID(invalidUserID)

	// Assert
	assert.Equal(t, entities.UserEntity{}, user)
}

func TestUserRepositoryGetByValidEmailReturnsUser(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	testUsers := setupUsers(userRepo)
	email := testUsers[0].Email

	// Act
	user := userRepo.GetByEmail(email)

	// Assert
	assert.Equal(t, testUsers[0], user)
}

func TestUserRepsitoryGetByInvalidEmailReturnsEmptyUser(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	invalidEmail := "test@email.nl"

	// Act
	user := userRepo.GetByEmail(invalidEmail)

	// Assert
	assert.Equal(t, entities.UserEntity{}, user)
}

func TestUserRepositoryCreateUserReturnsNewUser(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	testUsers := setupUsers(userRepo)
	userEntity := entities.UserEntity{
		ID:          3,
		FullName:    "Jil Doe",
		Email:       "jil@doe.org",
		AccountType: 0,
		Password:    "$2a$12$7/NMoWfjAzIZhkK/6S4yy.Pnvo1YbF1lxIR2ehNgKfz0xzVyExzZO",
		CreatedAt:   time.Date(2025, time.March, 31, 10, 30, 0, 0, time.UTC),
	}

	// Act
	user := userRepo.Create(userEntity)
	users := userRepo.GetAll()

	// Assert
	assert.Len(t, users, len(testUsers)+1)
	assert.Equal(t, userEntity, user)
}

func TestDeleteByValidIDReturnsTrue(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	testUsers := setupUsers(userRepo)
	userID := 1

	// Act
	isDeleted := userRepo.DeleteByID(userID)
	users := userRepo.GetAll()

	// Assert
	assert.Len(t, users, len(testUsers)-1)
	assert.True(t, isDeleted)
}

func TestDeleteByInvalidIDReturnsFalse(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	testUsers := setupUsers(userRepo)
	invalidUserID := 999

	// Act
	isDeleted := userRepo.DeleteByID(invalidUserID)
	users := userRepo.GetAll()

	// Assert
	assert.Len(t, users, len(testUsers))
	assert.False(t, isDeleted)
}

func TestUpdateValidUserReturnsUpdatedUser(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	testUsers := setupUsers(userRepo)
	// Update all user fields
	updatedUser := entities.UserEntity{
		ID:          1,
		FullName:    "Jonathan Doe",
		Email:       "jonathan@doe.it",
		AccountType: 1,
		Password:    "$2a$12$7/NMoWfjAzIZhkK/6S4yy.Pnvo1YbF1lxIR2ehNgKfz0xzVyExzZO",
		CreatedAt:   time.Date(2025, time.March, 31, 10, 30, 15, 0, time.UTC),
	}

	// Act
	user := userRepo.Update(updatedUser)

	// Assert
	assert.Equal(t, updatedUser, user)
	assert.NotNil(t, testUsers)
}

func TestSaveLastLoginTimeReturnsNil(t *testing.T) {
	// Arrange
	userRepo := NewTestUserRepository()
	mockUserID := setupUsers(userRepo)[0].ID

	// Act & Assert
	userRepo.SaveLastLoginTime(mockUserID)
}
