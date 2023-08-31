package apiRoutes

import (
	apiAuthRoutes "auth-api-cognito/internal/api/auth"
	cognito "auth-api-cognito/internal/auth"
	validatorUtil "auth-api-cognito/internal/utils/validator"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(public *gin.RouterGroup, private *gin.RouterGroup, v *validatorUtil.Validator, c *cognito.Cognito) {
	apiAuthRoutes.SetupRoutes(public, private, v, c)
}
