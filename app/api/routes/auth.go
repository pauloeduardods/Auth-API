package routes

import (
	"auth-api-cognito/internal/domain/handler"
)

func (r *Routes) SetupAuthRoutes() {
	authGroup := r.gin.Group("/auth")
	handler := handler.NewAuthHandler(r.services, r.utils)

	authGroup.POST("/login", handler.Login())
	authGroup.POST("/register", handler.Register())
	authGroup.POST("/confirm", handler.ConfirmSignUp())
	authGroup.GET("/info", r.middleware.AuthMiddleware(), handler.GetUser())
}
