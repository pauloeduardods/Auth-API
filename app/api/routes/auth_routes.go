package routes

import (
	handler "auth-api-cognito/api/handlers"
	"auth-api-cognito/api/middleware"

	"github.com/gin-gonic/gin"
)

func ConfigAuthRoutes(g *gin.Engine, m middleware.IAuthMiddleware, h handler.IAuthHandler) {
	authGroup := g.Group("/api/v1/auth")

	authGroup.POST("/login", h.Login())
	authGroup.POST("/register", h.Register())
	authGroup.POST("/confirm", h.ConfirmSignUp())
	authGroup.GET("/info", m.AuthMiddleware(), h.GetUser())
}
