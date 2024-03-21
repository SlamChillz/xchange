package token

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	CustomerID int32 `json:"customer_id"`
	jwt.RegisteredClaims `json:"claims"`
}

func (p *Payload) MarshalBinary() ([]byte, error) {
    return json.Marshal(p)
}

func NewPayload(customerId int32, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		CustomerID: customerId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: tokenId.String(),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	return payload, err
}
