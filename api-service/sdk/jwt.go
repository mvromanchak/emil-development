package sdk

import (
	"context"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const JwtDeviseIDField = "devise_id"

const bearerPrefix = "bearer"

type JWTSecured interface {
	Claims(ctx context.Context, secret string, method jwt.SigningMethod) (jwt.MapClaims, error)
	SetToken(string)
	Token() string
}

type JWTAuthorizer interface {
	Authorize(claims jwt.MapClaims) error
	JWTSecured
}

type JWTConfig struct {
	Secret     string
	ExtendTime time.Duration
	Method     jwt.SigningMethod
}

type simpleJWTSecured struct {
	jwtToken string
}

func (r *simpleJWTSecured) SetToken(token string) {
	r.jwtToken = token
}

func (r simpleJWTSecured) Token() string {
	return r.jwtToken
}
func NewJWTSecured() JWTSecured {
	return &simpleJWTSecured{}
}

func (r simpleJWTSecured) Claims(ctx context.Context, secret string, method jwt.SigningMethod) (jwt.MapClaims, error) {
	value := r.jwtToken
	if v := ctx.Value("jwt-token"); v != nil {
		//checkin the correct types under the interface
		jwtKey, ok := v.(string)
		if !ok {
			return nil, errors.New("claims context type assertion error")
		}
		value = jwtKey
	}
	return parseJWTToken(value, secret, method)
}

func parseJWTToken(jwtKey, secret string, method jwt.SigningMethod) (jwt.MapClaims, error) {
	if jwtKey == "" {
		return nil, errors.New("jwtKey is empty")
	}
	parsedToken, err := jwt.Parse(jwtKey, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken.Claims.(jwt.MapClaims), nil
}

func ExtractTokenFromBearer(bearer string) string {
	authHeaderParts := strings.Split(bearer, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != bearerPrefix {
		return ""
	}
	return authHeaderParts[1]
}
