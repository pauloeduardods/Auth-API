package middleware

import (
	"auth-api-cognito/internal/utils"
	"auth-api-cognito/pkg/jwtToken"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type IAuthMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
}

type AuthMiddleware struct {
	JwtToken *jwtToken.JwtToken
}

func NewAuthMiddleware(j *jwtToken.JwtToken) IAuthMiddleware {
	return &AuthMiddleware{
		JwtToken: j,
	}
}

func (a *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.Error(utils.NewApiError(http.StatusUnauthorized, "Authorization token missing"))
			c.Abort()
			return
		}

		splitToken := strings.Split(token, " ")

		//TODO: Check this
		a.JwtToken.CacheJWK()

		jwtToken, err := a.JwtToken.ParseJWT(splitToken[1])

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
