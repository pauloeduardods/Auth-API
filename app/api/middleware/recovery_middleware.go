package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryHandler(log *zap.Logger) gin.RecoveryFunc {
	return func(c *gin.Context, err any) {
		c.Next()
		log.Error("Error", zap.Any("recovered", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
		c.Abort()
		return
	}
}
