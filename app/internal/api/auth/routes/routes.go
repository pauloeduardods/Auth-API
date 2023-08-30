package authRoutes

import (
	"auth-api-cognito/internal/api/auth/controllers"
	cognito "auth-api-cognito/internal/auth"
	validatorUtil "auth-api-cognito/internal/utils/validator"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, v *validatorUtil.Validator, c *cognito.Cognito) {
	authGroup := r.Group("/auth")
	authGroup.POST("/login", controllers.Login(v, c))
	authGroup.POST("/register", controllers.Register(v, c))
	authGroup.POST("/confirm", controllers.ConfirmSignUp(v, c))
	authGroup.GET("/user", controllers.GetUser(v, c))
}
