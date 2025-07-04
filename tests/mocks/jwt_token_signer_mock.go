package mock_repositories

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
)

type MockJwtTokenSigner struct {
	mock.Mock
}

func (m *MockJwtTokenSigner) SignToken(claims jwt.Claims) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}
