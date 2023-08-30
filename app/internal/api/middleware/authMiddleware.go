package middleware

import (
	"auth-api-cognito/internal/utils"
	cognitoJwtVerify "auth-api-cognito/pkg/cognito-jwt-verify"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(j *cognitoJwtVerify.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.Error(utils.NewApiError(http.StatusUnauthorized, "Authorization token missing"))
			c.Abort()
			return
		}

		splitToken := strings.Split(token, " ")

		j.CacheJWK()

		jwtToken, err := j.ParseJWT(splitToken[1])

		c.Set("jwtToken", jwtToken)
		c.Set("user", jwtToken.Claims)

		if err != nil {
			c.Error(utils.NewApiError(http.StatusUnauthorized, err.Error()))
			c.Abort()
			return
		}

		c.Next()
	}
}
