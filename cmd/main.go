package main

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/pauloeduardods/auth-rest-api/internal/config"
	"github.com/pauloeduardods/auth-rest-api/internal/routes"
)

var r *gin.Engine

func init() {
	r = gin.Default()
	r.Use(gin.Logger())
	routes.SetupRoutes(r)
}

func main() {
	r.Run(":" + strconv.Itoa(config.EnvConfigs.Port))
}
