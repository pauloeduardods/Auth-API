package routes

import (
	"auth-api-cognito/api/middleware"
	services "auth-api-cognito/internal/domain/service"
	"auth-api-cognito/internal/utils"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	services   *services.Services
	utils      *utils.Utils
	gin        *gin.Engine
	middleware *middleware.Middleware
}

func NewRoutes(g *gin.Engine, m *middleware.Middleware, s *services.Services, u *utils.Utils) *Routes {
	return &Routes{
		services:   s,
		utils:      u,
		gin:        g,
		middleware: m,
	}
}
