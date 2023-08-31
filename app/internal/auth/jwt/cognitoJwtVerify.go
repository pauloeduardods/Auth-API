package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Auth struct {
	jwk               *JWK
	jwkURL            string
	cognitoRegion     string
	cognitoUserPoolID string
	log               *zap.Logger
}

type Config struct {
	CognitoRegion     string
	CognitoUserPoolID string
	Log               *zap.Logger
}

type JWK struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	} `json:"keys"`
}

func NewAuth(config *Config) *Auth {
	a := &Auth{
		cognitoRegion:     config.CognitoRegion,
		cognitoUserPoolID: config.CognitoUserPoolID,
		log:               config.Log,
	}

	a.jwkURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", a.cognitoRegion, a.cognitoUserPoolID)

	return a
}

func (a *Auth) CacheJWK() error {
	req, err := http.NewRequest("GET", a.jwkURL, nil)
	if err != nil {
		a.log.Error("Error creating JWK request", zap.Error(err))
		return err
	}

	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		a.log.Error("Error getting JWK", zap.Error(err))
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Error("Error reading JWK response body", zap.Error(err))
		return err
	}

	jwk := new(JWK)
	err = json.Unmarshal(body, jwk)
	if err != nil {
		a.log.Error("Error unmarshalling JWK", zap.Error(err))
		return err
	}

	a.jwk = jwk
	return nil
}

func (a *Auth) ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		key, err := convertKey(a.jwk.Keys[1].E, a.jwk.Keys[1].N)
		return key, err
	})
	if err != nil {
		a.log.Error("Error parsing JWT", zap.Error(err))
		return token, err
	}

	return token, nil
}

func (a *Auth) JWK() *JWK {
	return a.jwk
}

func (a *Auth) JWKURL() string {
	return a.jwkURL
}

func convertKey(rawE, rawN string) (*rsa.PublicKey, error) {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		return nil, err
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		return nil, err
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey, nil
}
