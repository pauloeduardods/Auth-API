package middleware

import "auth-api-cognito/pkg/jwtToken"

type Middleware struct {
	JwtToken *jwtToken.JwtToken
}

func NewMiddleware(j *jwtToken.JwtToken) *Middleware {
	return &Middleware{
		JwtToken: j,
	}
}
