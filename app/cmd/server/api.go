package server

import (
	apiRoutes "auth-api-cognito/internal/api"
	"auth-api-cognito/internal/middleware"
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

func (s *Server) SetupMiddlewareAndRouteGroup() {
	s.gin.Use(gin.CustomRecovery(middleware.RecoveryHandler(s.log)))
	s.gin.Use(gin.Logger())
	s.gin.Use(middleware.ErrorHandler(s.log))
	private := s.gin.Group("/api/v1")
	public := s.gin.Group("/api/v1")

	private.Use(middleware.AuthMiddleware(s.jwtVerify))

	s.privateRoute = private
	s.publicRoute = public
}

func (s *Server) SetupRoutes() {
	static.SetupStaticFiles(s.gin)
	apiRoutes.SetupRoutes(s.publicRoute, s.privateRoute, s.validator, s.cognito)
}
