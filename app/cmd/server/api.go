package server

import (
	"auth-api-cognito/api/handlers"
	"auth-api-cognito/api/middleware"
	"auth-api-cognito/api/routes"
	userService "auth-api-cognito/internal/domain/user/service"
	"auth-api-cognito/static"

	"github.com/gin-gonic/gin"
)

func (s *Server) SetupCors() {
	cors := middleware.Cors{
		Origin:      "*",
		Methods:     "GET, POST, PUT, DELETE, OPTIONS",
		Headers:     "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Auth-Token, X-Requested-With",
		Credentials: false,
	}
	s.gin.Use(cors.CorsMiddleware())
}

func (s *Server) SetupMiddlewares() {
	s.gin.Use(gin.CustomRecovery(middleware.RecoveryHandler(s.log)))
	s.gin.Use(gin.Logger())
	s.gin.Use(middleware.ErrorHandler(s.log))

}

func (s *Server) SetupApi() {
	static.SetupStaticFiles(s.gin)

	//Services

	authService := userService.NewAuthService(s.cognitoClient, s.cognitoClientId)

	//Handlers

	authHandler := handlers.NewAuthHandler(authService, s.utils)

	//Middlewares

	authMiddleware := middleware.NewAuthMiddleware(s.jwtToken)

	//Routes
	routes.ConfigAuthRoutes(s.gin, authMiddleware, authHandler)

}
