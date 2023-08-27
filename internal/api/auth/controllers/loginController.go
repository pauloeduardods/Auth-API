package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pauloeduardods/auth-rest-api/internal/api/auth/service"
	"github.com/pauloeduardods/auth-rest-api/internal/shared/utils"
)

func Login(c *gin.Context) {
	var login service.LoginInput

	if err := c.ShouldBindJSON(&login); err != nil {
		c.Error(err)
		return
	}

	err := utils.Validate(&login)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := login.Login()
	if err != nil {
		c.Error(err)
		return
	} else {
		c.JSON(http.StatusOK, res)
	}
}
