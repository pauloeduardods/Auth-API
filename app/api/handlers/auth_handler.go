package handlers

import (
	userServices "auth-api-cognito/internal/domain/user/service"
	"auth-api-cognito/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	ConfirmSignUp() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
}

type AuthHandler struct {
	authService userServices.IAuthService
	utils       *utils.Utils
}

func NewAuthHandler(s userServices.IAuthService, u *utils.Utils) IAuthHandler {
	return &AuthHandler{
		authService: s,
		utils:       u,
	}
}

func (a *AuthHandler) ConfirmSignUp() gin.HandlerFunc {
	return func(g *gin.Context) {
		var confirmSignUp userServices.ConfirmSignUpInput
		if err := g.ShouldBindJSON(&confirmSignUp); err != nil {
			g.Error(err)
			return
		}

		err := a.utils.Validate(&confirmSignUp)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.authService.ConfirmSignUp(confirmSignUp)
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
		getUser := userServices.GetUserInput{
			AccessToken: strings.Split(accessToken, " ")[1],
		}

		err := a.utils.Validate(&getUser)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.authService.GetUser(getUser)
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
		var login userServices.LoginInput
		if err := g.ShouldBindJSON(&login); err != nil {
			g.Error(err)
			return
		}

		err := a.utils.Validate(&login)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.authService.Login(login)
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
		var signUp userServices.SignUpInput
		if err := g.ShouldBindJSON(&signUp); err != nil {
			g.Error(err)
			return
		}

		err := a.utils.Validate(&signUp)
		if err != nil {
			g.Error(err)
			return
		}

		res, err := a.authService.SignUp(signUp)
		if err != nil {
			g.Error(err)
			return
		} else {
			g.JSON(http.StatusOK, res)
		}
	}
}
