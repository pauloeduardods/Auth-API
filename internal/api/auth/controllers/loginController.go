package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pauloeduardods/auth-rest-api/internal/api/auth"
	"github.com/pauloeduardods/auth-rest-api/internal/shared/utils"
)

func Login(c *gin.Context) {
	var login auth.LoginInput

	if err := c.ShouldBindJSON(&login); err != nil {
		c.Error(utils.NewApiError(http.StatusBadRequest, "Invalid request body", err.Error()))
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
