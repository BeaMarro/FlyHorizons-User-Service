package interfaces

import "github.com/gin-gonic/gin"

type GatewayAuthMiddleware interface {
	GatewayAuthMiddleware() gin.HandlerFunc
}
