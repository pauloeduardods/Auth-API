package routes

import (
	"github.com/gin-gonic/gin"
	authRoutes "github.com/pauloeduardods/auth-rest-api/internal/api/auth/routes"
)

func SetupRoutes(r *gin.Engine) {
	authRoutes.SetupRoutes(r)
}
