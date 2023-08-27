package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pauloeduardods/auth-rest-api/internal/config"
	"github.com/pauloeduardods/auth-rest-api/internal/shared/logger"
	"github.com/pauloeduardods/auth-rest-api/internal/shared/utils"
	"go.uber.org/zap"
)

func ErrorHandler() gin.HandlerFunc {
	appEnv := config.EnvConfigs.AppEnv
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
				return
			default:
				logger.Error("Error", zap.Error(e))
				if appEnv == "development" {
					c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": e.Error()})
					return
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
			}
		}
	}
}
