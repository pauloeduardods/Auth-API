package middleware

import (
	"auth-api-cognito/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.Error(utils.NewApiError(http.StatusUnauthorized, "Authorization token missing"))
			c.Abort()
			return
		}

		splitToken := strings.Split(token, " ")

		m.JwtToken.CacheJWK()

		jwtToken, err := m.JwtToken.ParseJWT(splitToken[1])

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
