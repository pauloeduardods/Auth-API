package middleware

import (
	"net/http"

	"auth-api-cognito/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/maragudk/env"
	"go.uber.org/zap"
)

func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case *utils.ApiError:
				c.AbortWithStatusJSON(e.StatusCode, e)
				c.Abort()
			case validator.ValidationErrors:
				errMsg := make(map[string]string)
				for _, fieldErr := range e {
					errMsg[fieldErr.Field()] = fieldErr.Tag()
				}
				c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
					"message": "Validation Error",
					"errors":  errMsg,
				})
				c.Abort()
				return
			default:
				appEnv := env.GetStringOrDefault("APP_ENV", "development")
				log.Error("Error", zap.Error(e))
				if appEnv == "development" {
					c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": e.Error()})
					c.Abort()
					return
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
				c.Abort()
			}
		}
	}
}
