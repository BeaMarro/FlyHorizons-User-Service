package mock_repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockGatewayAuthMiddleware struct {
	mock.Mock
	Role string
	ID   int
}

// Constructor to initialize MockGatewayAuthMiddleware with a role
func NewMockGatewayAuthMiddleware(role string, id int) *MockGatewayAuthMiddleware {
	return &MockGatewayAuthMiddleware{
		Role: role,
		ID:   id,
	}
}

// Middleware function now uses the role from the struct
func (m *MockGatewayAuthMiddleware) GatewayAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", m.ID)
		c.Set("sub", m.ID)
		c.Set("role", m.Role)
		c.Set("email", "test@email.com")

		c.Next()
	}
}
