package cognito

import (
	"context"

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
	return c.Client.InitiateAuth(context.TODO(), input)
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
