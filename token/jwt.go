package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	mininumAllowedSecretKeySize = 32
	ErrTokenIsInvalid		   = errors.New("token is invalid")
)

type JWT struct {
	secretKey string
}

func NewJWT(secretKey string) (*JWT, error) {
	if len(secretKey) < mininumAllowedSecretKeySize {
		return nil, fmt.Errorf("invalid secret key size: key must be at least %d characters", mininumAllowedSecretKeySize)
	}
	return &JWT{secretKey}, nil
}

func (jwtoken *JWT) CreateToken(useremail string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(useremail, duration)
	if err != nil {
		return "", payload, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(jwtoken.secretKey))
	return tokenString, payload, err
}

func (jwtoken *JWT) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ! ok {
			return nil, fmt.Errorf("invalid token signing method: %v", token.Header["alg"])
		}
		return []byte(jwtoken.secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		// Log the error
		return nil, ErrTokenIsInvalid
	}
	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrTokenIsInvalid
	}
	return payload, nil
}
