package mock_repositories

import (
	"flyhorizons-userservice/models"
	"flyhorizons-userservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

var _ interfaces.UserService = (*MockUserService)(nil)

func (m *MockUserService) GetAll() []models.User {
	args := m.Called()
	return args.Get(0).([]models.User)
}

func (m *MockUserService) GetByID(userID int) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UserExists(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *MockUserService) Create(user models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteByID(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserService) Update(user models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
