package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pauloeduardods/auth-rest-api/internal/api/auth/controllers"
)

func SetupRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.POST("/login", controllers.Login)
}
