package server

import (
	authRoutes "auth-api-cognito/internal/api/auth/routes"
	"auth-api-cognito/internal/api/middleware"
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
	s.gin.Use(cors.Cors())
}

func (s *Server) SetupMiddleware() {
	s.gin.Use(middleware.ErrorHandler(s.log))
	s.gin.Use(gin.Recovery())
	s.gin.Use(gin.Logger())
}

func (s *Server) SetupRoutes() {
	static.SetupStaticFiles(s.gin)
	authRoutes.SetupRoutes(s.gin, s.validator, s.cognito)
}
