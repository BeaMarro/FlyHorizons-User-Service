package mock_repositories

import (
	entities "flyhorizons-userservice/repositories/entity"
	"flyhorizons-userservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

var _ interfaces.UserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) GetAll() []entities.UserEntity {
	args := m.Called()
	return args.Get(0).([]entities.UserEntity)
}

func (m *MockUserRepository) GetByID(userID int) entities.UserEntity {
	args := m.Called(userID)
	return args.Get(0).(entities.UserEntity)
}

func (m *MockUserRepository) GetByEmail(email string) entities.UserEntity {
	args := m.Called(email)
	return args.Get(0).(entities.UserEntity)
}

func (m *MockUserRepository) Create(user entities.UserEntity) entities.UserEntity {
	args := m.Called(user)
	return args.Get(0).(entities.UserEntity)
}

func (m *MockUserRepository) DeleteByID(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *MockUserRepository) Update(user entities.UserEntity) entities.UserEntity {
	args := m.Called(user)
	return args.Get(0).(entities.UserEntity)
}

func (m *MockUserRepository) SaveLastLoginTime(id int) {
	m.Called(id)
}
