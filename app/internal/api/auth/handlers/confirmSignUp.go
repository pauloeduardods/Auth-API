package handlers

import (
	"net/http"

	cognito "auth-api-cognito/internal/auth"
	validatorUtil "auth-api-cognito/internal/utils/validator"

	"github.com/gin-gonic/gin"
)

func ConfirmSignUp(v *validatorUtil.Validator, c *cognito.Cognito) gin.HandlerFunc {
	return func(g *gin.Context) {
		var confirmSignUp cognito.ConfirmSignUpInput
		if err := g.ShouldBindJSON(&confirmSignUp); err != nil {
			g.Error(err)
			return
		}

		err := v.Validate(&confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := c.ConfirmSignUp(confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}
