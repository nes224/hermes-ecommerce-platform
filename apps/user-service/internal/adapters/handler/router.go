package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	envString = "production"
)

func SetupRouter(env string, authHandler *AuthHandler) *gin.Engine {
	if env == envString {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "UP",
			"environment": env,
			"service":     "user-service",
		})
	})

	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/signup", authHandler.SignUp)
			authRoutes.POST("/signin", authHandler.SignIn)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
			authRoutes.POST("/logout", authHandler.Logout)
		}
	}

	return router
}
