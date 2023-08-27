package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	cognitoSvc *cognito.CognitoIdentityProvider
	userPoolID = os.Getenv("COGNITO_USER_POOL_ID")
	region     = os.Getenv("COGNITO_REGION")
	log        *zap.Logger
)

func handleErr(key string, err error, isCritical bool) {
	if err != nil {
		log.Error(key, zap.Error(err))
		if isCritical {
			panic(err)
		}
	}
}

var ErrUserNotFound = errors.New("user not found")

func init() {
	var cfg zap.Config

	cfg = zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	handleErr("CREATE_SESSION_ERROR", err, true)

	cognitoSvc = cognito.New(sess)
}

func userExists(_ *cognito.AdminGetUserOutput, err error) bool {
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return false
		}
		handleErr("USER_EXISTS_ERROR", err, true)
	}
	return true
}

func lookupUserByEmail(email string) (*cognito.AdminGetUserOutput, error) {
	input := &cognito.AdminGetUserInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(email),
	}
	output, err := cognitoSvc.AdminGetUser(input)
	log.Info("LOOKUP_USER", zap.Any("input", input), zap.Any("output", output), zap.Error(err))
	if err != nil {
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return nil, ErrUserNotFound
		}
		handleErr("LOOKUP_USER_ERROR", err, true)
	}
	return output, nil
}

func linkFederatedUserToExistingUser(event events.CognitoEventUserPoolsPreSignup) {
	splitUserName := strings.Split(event.UserName, "_")

	input := &cognito.AdminLinkProviderForUserInput{
		UserPoolId: &event.UserPoolID,
		DestinationUser: &cognito.ProviderUserIdentifierType{
			ProviderAttributeValue: aws.String(event.Request.UserAttributes["email"]),
			ProviderName:           aws.String("Cognito"),
		},
		SourceUser: &cognito.ProviderUserIdentifierType{
			ProviderAttributeName:  aws.String("Cognito_Subject"),
			ProviderAttributeValue: aws.String(splitUserName[1]),
			ProviderName:           aws.String(splitUserName[0]),
		},
	}
	link, err := cognitoSvc.AdminLinkProviderForUser(input)
	handleErr("LINK_USER_ERROR", err, true)
	log.Info("LINK_USER", zap.Any("input", input), zap.Any("output", link))
}

func createNewUser(event events.CognitoEventUserPoolsPreSignup) *cognito.AdminCreateUserOutput {
	log.Info("CREATE_USER", zap.Any("event", event))
	email := event.Request.UserAttributes["email"]
	input := &cognito.AdminCreateUserInput{
		UserPoolId:    &event.UserPoolID,
		Username:      &email,
		MessageAction: aws.String("SUPPRESS"),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
		},
	}
	user, err := cognitoSvc.AdminCreateUser(input)
	log.Info("CREATE_USER", zap.Any("input", input), zap.Any("output", user), zap.Error(err))
	handleErr("CREATE_USER_ERROR", err, true)
	return user
}

func handler(ctx context.Context, event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	log.Info("START_EVENT", zap.Any("event", event))
	if event.TriggerSource != "PreSignUp_ExternalProvider" {
		return event, nil
	}

	email := event.Request.UserAttributes["email"]
	user, err := lookupUserByEmail(email)
	userExistsInCurrentPool := userExists(user, err)

	if !userExistsInCurrentPool {
		createNewUser(event)
	}

	linkFederatedUserToExistingUser(event)

	log.Info("END_EVENT", zap.Any("event", event))

	return event, nil
}

func main() {
	lambda.Start(handler)
}
