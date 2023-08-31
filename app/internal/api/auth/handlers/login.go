package handlers

import (
	"net/http"

	cognito "auth-api-cognito/internal/auth"
	validatorUtil "auth-api-cognito/internal/utils/validator"

	"github.com/gin-gonic/gin"
)

func Login(v *validatorUtil.Validator, c *cognito.Cognito) gin.HandlerFunc {
	return func(g *gin.Context) {
		var login cognito.LoginInput
		if err := g.ShouldBindJSON(&login); err != nil {
			g.Error(err)
			return
		}

		err := v.Validate(&login)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := c.Login(login)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res.AuthenticationResult)
		}
	}
}
