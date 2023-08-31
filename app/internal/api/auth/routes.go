package apiAuthRoutes

import (
	"auth-api-cognito/internal/api/auth/handlers"
	cognito "auth-api-cognito/internal/auth"
	validatorUtil "auth-api-cognito/internal/utils/validator"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(public *gin.RouterGroup, private *gin.RouterGroup, v *validatorUtil.Validator, c *cognito.Cognito) {
	publicGroup := public.Group("/auth")
	privateGroup := private.Group("/auth")

	publicGroup.POST("/login", handlers.Login(v, c))
	publicGroup.POST("/register", handlers.Register(v, c))
	publicGroup.POST("/confirm", handlers.ConfirmSignUp(v, c))
	privateGroup.GET("/info", handlers.GetUser(v, c))
}
