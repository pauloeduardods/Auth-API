package handlers

import (
	"net/http"
	"strings"

	cognito "auth-api-cognito/internal/auth"
	validatorUtil "auth-api-cognito/internal/utils/validator"

	"github.com/gin-gonic/gin"
)

func GetUser(v *validatorUtil.Validator, c *cognito.Cognito) gin.HandlerFunc {
	return func(g *gin.Context) {
		accessToken := g.GetHeader("Authorization")
		getUser := cognito.GetUserInput{
			AccessToken: strings.Split(accessToken, " ")[1],
		}

		err := v.Validate(&getUser)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := c.GetUser(getUser)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}
