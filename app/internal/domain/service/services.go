package services

import (
	cognitoService "auth-api-cognito/internal/domain/service/cognito"
)

type Services struct {
	Cognito *cognitoService.CognitoService
}

func NewServices(c *cognitoService.CognitoService) *Services {
	return &Services{
		Cognito: c,
	}
}
