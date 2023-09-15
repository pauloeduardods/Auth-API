package api

import (
	"auth-api-cognito/api/middleware"
	"auth-api-cognito/api/routes"
	services "auth-api-cognito/internal/domain/service"
	"auth-api-cognito/internal/utils"
	"auth-api-cognito/pkg/jwtToken"

	"github.com/gin-gonic/gin"
)

type API struct {
	Router   *gin.Engine
	Services *services.Services
	Utils    *utils.Utils
	JwtToken *jwtToken.JwtToken
}

func NewAPI(r *gin.Engine, s *services.Services, u *utils.Utils, j *jwtToken.JwtToken) *API {
	return &API{
		Router:   r,
		Services: s,
		Utils:    u,
		JwtToken: j,
	}
}

func (a *API) SetupApi() {
	middleware := middleware.NewMiddleware(a.JwtToken)
	routes := routes.NewRoutes(a.Router, middleware, a.Services, a.Utils)
	routes.SetupAuthRoutes()
}
