package routes

import (
	"flyhorizons-userservice/models/request"
	"flyhorizons-userservice/services/interfaces"
	"flyhorizons-userservice/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, loginService interfaces.LoginService) {
	router.POST("/login", func(c *gin.Context) {
		var loginRequest request.LoginRequest

		// Bind the JSON request body to the loginRequest struct
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the client IP address
		ipAddress := utils.GetIPAddress(c.Request)

		// Call the login service with email, password, and IP
		loginResponse, err := loginService.Login(loginRequest, ipAddress)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, loginResponse)
	})
}
