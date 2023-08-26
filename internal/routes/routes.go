package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pauloeduardods/auth-rest-api/internal/api/auth/routes"
)

func SetupRoutes(r *gin.Engine) {
	routes.SetupAuthRoutes(r)
}
