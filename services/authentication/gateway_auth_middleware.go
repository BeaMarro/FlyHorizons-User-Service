package authentication

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type GatewayAuthMiddlewareHandler struct{}

func (g *GatewayAuthMiddlewareHandler) GatewayAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Load .env file
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, relying on environment variables")
		}

		JwtSecret := []byte(os.Getenv("JWT_SECRET"))

		// Get JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the JWT token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtSecret, nil
		})

		if err != nil || !token.Valid {
			fmt.Println("JWT parsing failed:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid JWT token"})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid JWT claims"})
			return
		}

		// TODO: Delete the logging after implementing the refreshing authentication tokens
		fmt.Println("--- Extracted JWT Claims ---")
		for k, v := range claims {
			fmt.Printf("%s: %v\n", k, v)
		}

		// Set claims
		if sub, ok := claims["sub"].(float64); ok {
			c.Set("user_id", int(sub))
			c.Set("sub", int(sub))
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}
		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}

		c.Next()
	}
}

func NewGatewayAuthMiddleware() *GatewayAuthMiddlewareHandler {
	return &GatewayAuthMiddlewareHandler{}
}
