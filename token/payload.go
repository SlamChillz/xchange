package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
}

func NewPayload(userEmail string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: tokenId.String(),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	return payload, err
}
