package service

import (
	"auth-api-cognito/internal/utils"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type AuthService struct {
	Client   *cognito.Client
	clientId string
}

type IAuthService interface {
	Login(l LoginInput) (*cognito.InitiateAuthOutput, error)
	SignUp(s SignUpInput) (*cognito.SignUpOutput, error)
	UserInformation(accessToken string) (*cognito.GetUserOutput, error)
	ConfirmSignUp(s ConfirmSignUpInput) (*cognito.ConfirmSignUpOutput, error)
	GetUser(g GetUserInput) (*cognito.GetUserOutput, error)
}

type LoginInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

type SignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
	Name     string `json:"name" binding:"required" validate:"min=3,max=50"`
}

type ConfirmSignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Code     string `json:"code" binding:"required" validate:"numeric"`
}

type GetUserInput struct {
	AccessToken string `json:"accessToken" binding:"required"`
}

func NewAuthService(c *cognito.Client, clientId string) *AuthService {
	return &AuthService{
		Client:   c,
		clientId: clientId,
	}
}

func (c *AuthService) Login(l LoginInput) (*cognito.InitiateAuthOutput, error) {
	input := &cognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": l.Username,
			"PASSWORD": l.Password,
		},
		ClientId: aws.String(c.clientId),
	}
	out, err := c.Client.InitiateAuth(context.TODO(), input)
	if err != nil {
		errorType := err.Error()
		if strings.Contains(errorType, "NotAuthorizedException") {
			return nil, utils.NewApiError(401, "Invalid username or password")
		}
		if strings.Contains(errorType, "PasswordResetRequiredException") {
			return nil, utils.NewApiError(401, "Password reset required")
		}
		if strings.Contains(errorType, "UserNotConfirmedException") {
			return nil, utils.NewApiError(401, "User not confirmed")
		}
		return nil, err
	}
	return out, nil
}

func (c *AuthService) SignUp(s SignUpInput) (*cognito.SignUpOutput, error) {
	input := &cognito.SignUpInput{
		ClientId: aws.String(c.clientId),
		Username: aws.String(s.Username),
		Password: aws.String(s.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(s.Name),
			},
		},
	}
	return c.Client.SignUp(context.TODO(), input)
}

func (c *AuthService) UserInformation(accessToken string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	return c.Client.GetUser(context.TODO(), input)
}

func (c *AuthService) ConfirmSignUp(s ConfirmSignUpInput) (*cognito.ConfirmSignUpOutput, error) {
	input := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientId),
		Username:         aws.String(s.Username),
		ConfirmationCode: aws.String(s.Code),
	}
	return c.Client.ConfirmSignUp(context.Background(), input)
}

func (c *AuthService) GetUser(g GetUserInput) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: &g.AccessToken,
	}
	return c.Client.GetUser(context.Background(), input)
}
