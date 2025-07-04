package mock_repositories

import (
	"flyhorizons-userservice/models/request"
	"flyhorizons-userservice/models/response"
	"flyhorizons-userservice/services/interfaces"

	"github.com/stretchr/testify/mock"
)

type MockLoginService struct {
	mock.Mock
}

var _ interfaces.LoginService = (*MockLoginService)(nil)

func (m *MockLoginService) Login(loginRequest request.LoginRequest) (*response.LoginResponse, error) {
	args := m.Called(loginRequest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.LoginResponse), args.Error(1)
}
