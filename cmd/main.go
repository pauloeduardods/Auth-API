package main

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/pauloeduardods/auth-rest-api/internal/config"
	"github.com/pauloeduardods/auth-rest-api/internal/routes"
	"github.com/pauloeduardods/auth-rest-api/internal/shared/middleware"
)

var (
	r    *gin.Engine
	cors = middleware.Cors{
		Origin:      "https://dev.meuguru.net",
		Methods:     "GET, POST, PUT, DELETE, OPTIONS",
		Headers:     "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Auth-Token, X-Requested-With",
		Credentials: false,
	}
)

func init() {
	r = gin.Default()
	r.Use(cors.Cors())
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(gin.Recovery())
	routes.SetupRoutes(r)
}

func main() {
	r.Run(":" + strconv.Itoa(config.EnvConfigs.Port))
}
