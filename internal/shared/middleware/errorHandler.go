package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
