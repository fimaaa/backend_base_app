package middleware

import (
	"fmt"
	"reflect"

	"github.com/dgrijalva/jwt-go"
)

type (
	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		// Signing key to validate token.
		// Required.
		SigningKey interface{}

		// Signing method, used to check token signing method.
		// Optional. Default value HS256.
		SigningMethod string

		// Context key to store user information from the token into context.
		// Optional. Default value "user".
		ContextKey string

		// Claims are extendable claims data defining token content.
		// Optional. Default value jwt.MapClaims
		Claims jwt.Claims

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup string

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string

		keyFunc jwt.Keyfunc
	}
)

// Algorithms
const (
	AlgorithmHS256 = "HS256"
)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = JWTConfig{
		SigningMethod: AlgorithmHS256,
		ContextKey:    "user",
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
	}
)

func GetMapClaimByKeyJwtToken(secret, auth, key string) (interface{}, error) {
	var err error

	config := JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(secret),
	}

	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}

	config.keyFunc = func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}
		return config.SigningKey, nil
	}

	token := new(jwt.Token)

	// Issue #647, #656
	if _, ok := config.Claims.(jwt.MapClaims); ok {
		token, err = jwt.Parse(auth, config.keyFunc)
	} else {
		claims := reflect.ValueOf(config.Claims).Interface().(jwt.Claims)
		token, err = jwt.ParseWithClaims(auth, claims, config.keyFunc)
	}
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims[key], nil
}
