package config

import (
	"fmt"

	"github.com/maragudk/env"
)

type AppConfig struct {
	Host              string
	Port              int
	CognitoClientId   string
	CognitoRegion     string
	CognitoUserPoolID string
	AppEnv            string
}

var release string

func LoadConfig() (*AppConfig, error) {
	_ = env.Load(".env")

	appEnv := env.GetStringOrDefault("APP_ENV", "development")
	host := env.GetStringOrDefault("HOST", "localhost")
	port := env.GetIntOrDefault("PORT", 4000)
	cognitoClientId := env.GetStringOrDefault("COGNITO_CLIENT_ID", "")
	cognitoUserPoolID := env.GetStringOrDefault("COGNITO_USER_POOL_ID", "")
	cognitoRegion := env.GetStringOrDefault("COGNITO_REGION", "")

	if host == "" {
		return nil, fmt.Errorf("Host is not defined")
	}
	if cognitoClientId == "" {
		return nil, fmt.Errorf("CognitoClientId is not defined")
	}
	if cognitoUserPoolID == "" {
		return nil, fmt.Errorf("CognitoUserPoolID is not defined")
	}
	if cognitoRegion == "" {
		return nil, fmt.Errorf("CognitoRegion is not defined")
	}

	config := &AppConfig{
		Host:              host,
		Port:              port,
		CognitoClientId:   cognitoClientId,
		CognitoRegion:     cognitoRegion,
		CognitoUserPoolID: cognitoUserPoolID,
		AppEnv:            appEnv,
	}
	return config, nil
}
