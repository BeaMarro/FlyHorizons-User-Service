package main

import (
	"flyhorizons-userservice/config"
	"flyhorizons-userservice/internal/health"
	"flyhorizons-userservice/internal/metrics"
	"flyhorizons-userservice/repositories"
	"flyhorizons-userservice/routes"
	"flyhorizons-userservice/services"
	"flyhorizons-userservice/services/authentication"
	"flyhorizons-userservice/services/converter"
	"flyhorizons-userservice/services/validation"

	"github.com/gin-gonic/gin"

	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	healthcfg "github.com/tavsec/gin-healthcheck/config"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	// Initialize RabbitMQ for messaging
	config.InitializeRabbitMQ()
	defer config.RabbitMQClient.Connection.Close()
	defer config.RabbitMQClient.Connection.Channel()

	router := gin.Default()

	// Initialize repository
	baseRepo := &repositories.BaseRepository{}
	dbCheck := health.DatabaseCheck{Repository: baseRepo}

	// --- Health checks setup ---
	conf := healthcfg.DefaultConfig()
	rabbitMQCheck := health.RabbitMQCheck{}
	healthcheck.New(router, conf, []checks.Check{dbCheck, rabbitMQCheck})

	// --- Metrics setup ---
	metrics.RegisterMetricsRoutes(router, dbCheck, rabbitMQCheck)

	// --- Microservice setup ---
	userRepo := repositories.NewUserRepository(baseRepo)

	// Initialize services
	userConverter := converter.UserConverter{}
	passwordValidator := validation.PasswordValidator{}
	accountHashing := authentication.NewAccountHashing()
	jwtSigner := authentication.NewJwtTokenSigner()
	oauthSigner := services.NewOAuthTokenSigner(jwtSigner)

	// Authentication middlware
	gatewayAuthMiddleware := authentication.NewGatewayAuthMiddleware()
	loginService := services.NewLoginService(userRepo, userConverter, oauthSigner)
	userService := services.NewUserService(userRepo, accountHashing, passwordValidator, userConverter)

	// Register routes
	routes.RegisterUserRoutes(router, userService, gatewayAuthMiddleware)
	routes.RegisterAuthRoutes(router, loginService)

	// Run the microservice
	router.Run(":8081")
}
