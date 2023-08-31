package handlers

import (
	"net/http"

	cognito "auth-api-cognito/internal/auth"
	validatorUtil "auth-api-cognito/internal/utils/validator"

	"github.com/gin-gonic/gin"
)

func Register(v *validatorUtil.Validator, c *cognito.Cognito) gin.HandlerFunc {
	return func(g *gin.Context) {
		var signUp cognito.SignUpInput
		if err := g.ShouldBindJSON(&signUp); err != nil {
			g.Error(err)
			return
		}

		err := v.Validate(&signUp)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := c.SignUp(signUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}
