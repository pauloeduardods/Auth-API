package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pauloeduardods/auth-rest-api/internal/api/auth"
)

func Login(c *gin.Context) {
	var login auth.Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := login.Login()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, res)
	}
}
