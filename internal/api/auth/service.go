package auth

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pauloeduardods/auth-rest-api/internal/config"
	"github.com/pauloeduardods/auth-rest-api/internal/utils"
)

var (
	cognitoClient *cognito.CognitoIdentityProvider
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.EnvConfigs.CognitoRegion),
	}))
	cognitoClient = cognito.New(sess)
}

type Login struct {
	Username string `json:"username" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=8"`
}

func (l *Login) Login() (*cognito.InitiateAuthOutput, error) {
	err := utils.Validate(l)

	if err != nil {
		return nil, err
	}

	input := &cognito.InitiateAuthInput{
		AuthFlow: aws.String(cognito.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(l.Username),
			"PASSWORD": aws.String(l.Password),
		},
		ClientId: aws.String(config.EnvConfigs.CognitoClientId),
	}
	return cognitoClient.InitiateAuth(input)
}

func SignUp(username, password string) (*cognito.SignUpOutput, error) {
	input := &cognito.SignUpInput{
		ClientId: aws.String(config.EnvConfigs.CognitoClientId),
		Username: aws.String(username),
		Password: aws.String(password),
	}
	return cognitoClient.SignUp(input)
}