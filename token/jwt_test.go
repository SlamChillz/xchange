package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func TestValidJWT(t *testing.T) {
	jwt, err := NewJWT(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwt)
	userEmail := utils.RandomEmail()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	tokenString, payload, err := jwt.CreateToken(userEmail, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotEmpty(t, payload)

	payload, err = jwt.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, userEmail, payload.UserEmail)
	require.WithinDuration(t, issuedAt, payload.RegisteredClaims.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiredAt, payload.RegisteredClaims.ExpiresAt.Time, time.Second)
}

func TestInvalidJWT(t *testing.T) {
	payload, err := NewPayload(utils.RandomEmail(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	jwt, err := NewJWT(utils.RandomString(32))
	require.NoError(t, err)

	payload, err = jwt.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrTokenIsInvalid.Error())
	require.Nil(t, payload)
}

func TestExpiredJWT(t *testing.T) {
	jwtToken, err := NewJWT(utils.RandomString(32))
	require.NoError(t, err)

	tokenString, payload, err := jwtToken.CreateToken(utils.RandomEmail(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotEmpty(t, payload)

	payload, err = jwtToken.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, jwt.ErrTokenExpired.Error())
	require.Nil(t, payload)
}
