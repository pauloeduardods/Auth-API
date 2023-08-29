package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/pauloeduardods/auth-rest-api/internal/config"
	"github.com/pauloeduardods/auth-rest-api/internal/shared/utils"
)

var (
	cognitoClient *cognito.Client
)

func init() {
	cfg := config.GetAWSConfig()
	cognitoClient = cognito.NewFromConfig(cfg)
}

type LoginInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

func (l *LoginInput) Login() (*cognito.InitiateAuthOutput, error) {
	input := &cognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": l.Username,
			"PASSWORD": l.Password,
		},
		ClientId: aws.String(config.EnvConfigs.CognitoClientId),
	}
	return cognitoClient.InitiateAuth(context.TODO(), input)
}

type SignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

func (s *SignUpInput) SignUp() (*cognito.SignUpOutput, error) {
	err := utils.Validate(s)

	if err != nil {
		return nil, err
	}
	input := &cognito.SignUpInput{
		ClientId: aws.String(config.EnvConfigs.CognitoClientId),
		Username: aws.String(s.Username),
		Password: aws.String(s.Password),
	}
	return cognitoClient.SignUp(context.TODO(), input)
}

func ValidateToken(accessToken string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	return cognitoClient.GetUser(context.TODO(), input)
}

func UserInformation(accessToken string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	return cognitoClient.GetUser(context.TODO(), input)
}
