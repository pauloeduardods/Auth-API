package middleware

import (
	"auth-api-cognito/internal/auth/jwt"
	"auth-api-cognito/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(j *jwt.Auth) gin.HandlerFunc {
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

		if err != nil {
			c.Error(utils.NewApiError(http.StatusUnauthorized, err.Error()))
			c.Abort()
			return
		}

		c.Set("jwtToken", jwtToken)
		c.Set("user", jwtToken.Claims)

		c.Next()
	}
}
