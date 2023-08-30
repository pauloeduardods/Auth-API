package cognito

import (
	"auth-api-cognito/internal/utils"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type Cognito struct {
	Client   *cognito.Client
	clientId string
}

type Options struct {
	AWSConfig aws.Config
	ClientId  string
}

func New(opts Options) *Cognito {
	return &Cognito{
		Client:   cognito.NewFromConfig(opts.AWSConfig),
		clientId: opts.ClientId,
	}
}

type LoginInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

func (c *Cognito) Login(l LoginInput) (*cognito.InitiateAuthOutput, error) {
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

type SignUpInput struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

func (c *Cognito) SignUp(s SignUpInput) (*cognito.SignUpOutput, error) {
	// err := c.validator.Struct(s)

	// if err != nil {
	// 	return nil, err
	// }
	input := &cognito.SignUpInput{
		ClientId: aws.String(c.clientId),
		Username: aws.String(s.Username),
		Password: aws.String(s.Password),
	}
	return c.Client.SignUp(context.TODO(), input)
}

func (c *Cognito) UserInformation(accessToken string) (*cognito.GetUserOutput, error) {
	input := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	return c.Client.GetUser(context.TODO(), input)
}
