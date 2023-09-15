package handler

import (
	"net/http"
	"strings"

	services "auth-api-cognito/internal/domain/service"
	cognitoService "auth-api-cognito/internal/domain/service/cognito"
	"auth-api-cognito/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	services *services.Services
	utils    *utils.Utils
}

func NewAuthHandler(s *services.Services, u *utils.Utils) *AuthHandler {
	return &AuthHandler{
		services: s,
		utils:    u,
	}
}

func (a *AuthHandler) ConfirmSignUp() gin.HandlerFunc {
	return func(g *gin.Context) {
		var confirmSignUp cognitoService.ConfirmSignUpInput
		if err := g.ShouldBindJSON(&confirmSignUp); err != nil {
			g.Error(err)
			return
		}

		err := a.utils.Validate(&confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.services.Cognito.ConfirmSignUp(confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}

func (a *AuthHandler) GetUser() gin.HandlerFunc {
	return func(g *gin.Context) {
		accessToken := g.GetHeader("Authorization")
		getUser := cognitoService.GetUserInput{
			AccessToken: strings.Split(accessToken, " ")[1],
		}

		err := a.utils.Validate(&getUser)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.services.Cognito.GetUser(getUser)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}

func (a *AuthHandler) Login() gin.HandlerFunc {
	return func(g *gin.Context) {
		var login cognitoService.LoginInput
		if err := g.ShouldBindJSON(&login); err != nil {
			g.Error(err)
			return
		}

		err := a.utils.Validate(&login)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.services.Cognito.Login(login)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res.AuthenticationResult)
		}
	}
}

func (a *AuthHandler) Register() gin.HandlerFunc {
	return func(g *gin.Context) {
		var signUp cognitoService.SignUpInput
		if err := g.ShouldBindJSON(&signUp); err != nil {
			g.Error(err)
			return
		}

		err := a.utils.Validate(&signUp)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.services.Cognito.SignUp(signUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}
