package routes

import (
	"flyhorizons-userservice/models"
	"flyhorizons-userservice/services/errors"
	"flyhorizons-userservice/services/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userService interfaces.UserService, authMiddleware interfaces.GatewayAuthMiddleware) {
	// Public route
	router.POST("/users", func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		postUser, err := userService.Create(user)
		if err != nil {
			if _, ok := err.(*errors.UserExistsError); ok {
				ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
				return
			}
			if _, ok := err.(*errors.InsufficientPasswordLengthError); ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if _, ok := err.(*errors.InsufficientPasswordComplexityError); ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if _, ok := err.(*errors.CommonPasswordError); ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, postUser)
	})

	userGroup := router.Group("/users")
	userGroup.Use(authMiddleware.GatewayAuthMiddleware())

	// Protected routes
	// Only accessible by admins
	userGroup.GET("/", func(ctx *gin.Context) {
		role, exists := ctx.Get("role")

		if !exists || role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: admin access required"})
			return
		}

		users := userService.GetAll()
		ctx.JSON(http.StatusOK, users)
	})

	// Only accessible by users with the matching ID
	userGroup.GET("/:userID", func(ctx *gin.Context) {
		userIDString := ctx.Param("userID")
		userID, err := strconv.Atoi(userIDString)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
			return
		}

		tokenUserID := ctx.GetInt("user_id")

		if tokenUserID != userID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot access the account belonging to another user"})
			return
		}

		user, err := userService.GetByID(userID)
		if err != nil {
			if _, ok := err.(*errors.UserNotFoundError); ok {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	})

	// Only accessible by users with the matching ID
	userGroup.DELETE("/:ID", func(ctx *gin.Context) {
		userIDString := ctx.Param("ID")
		userID, err := strconv.Atoi(userIDString)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
			return
		}

		tokenUserID := ctx.GetInt("sub")

		if tokenUserID != userID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot delete the account of a different user"})
			return
		}

		success, err := userService.DeleteByID(userID)
		if err != nil {
			if _, ok := err.(*errors.UserNotFoundError); ok {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if success {
			ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete user"})
		}
	})

	// Only accessible by users with the matching ID
	userGroup.PUT("/", func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tokenUserID := ctx.GetInt("sub")

		if user.ID != tokenUserID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: cannot update the account belonging to another user"})
			return
		}

		putUser, err := userService.Update(user)
		if err != nil {
			if _, ok := err.(*errors.UserNotFoundError); ok {
				ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			if _, ok := err.(*errors.InsufficientPasswordLengthError); ok { // Other 2 types of exceptions to throw this type of custom exception
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if _, ok := err.(*errors.InsufficientPasswordComplexityError); ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if _, ok := err.(*errors.CommonPasswordError); ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, putUser)
	})
}
