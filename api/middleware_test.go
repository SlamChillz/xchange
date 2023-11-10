package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slamchillz/xchange/token"
	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func addAccessTokenToRequestHeader(
	t *testing.T,
	req *http.Request,
	token *token.JWT,
	authenticationScheme string,
	customerId int32,
	duration time.Duration,
) {
	tokenString, payload, err := token.CreateToken(customerId, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	authenticationHeader := fmt.Sprintf("%s %s", authenticationScheme, tokenString)
	req.Header.Set(AUTHENTICATIONHEADER, authenticationHeader)
}

func TestAuthenticationMiddleware(t *testing.T) {
	customer := utils.RandomNumber()
	testCases := []struct {
		name string
		setUpAuth func (t *testing.T, req *http.Request, authToken *token.JWT)
		responseCheck func (t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setUpAuth: func (t *testing.T, req *http.Request, authToken *token.JWT) {
				addAccessTokenToRequestHeader(t, req, authToken, AUTHENTICATIONSCHEME, customer, time.Minute)
			},
			responseCheck: func (t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "UnrecognizedAuthenticationHeaderValue",
			setUpAuth: func (t *testing.T, req *http.Request, authToken *token.JWT) {
				req.Header.Set(AUTHENTICATIONHEADER, "Bearer valid invalid")
			},
			responseCheck: func (t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "NoAuthenticationeader",
			setUpAuth: func (t *testing.T, req *http.Request, authToken *token.JWT) {},
			responseCheck: func (t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "UnsupportedAuthenticationHeader",
			setUpAuth: func (t *testing.T, req *http.Request, authToken *token.JWT) {
				addAccessTokenToRequestHeader(t, req, authToken, "Basic", customer, time.Minute)
			},
			responseCheck: func (t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "InvalidAuthenticationHeader",
			setUpAuth: func (t *testing.T, req *http.Request, authToken *token.JWT) {
				req.Header.Set(AUTHENTICATIONHEADER, "Bearer invalid")
			},
			responseCheck: func (t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setUpAuth: func (t *testing.T, req *http.Request, authToken *token.JWT) {
				addAccessTokenToRequestHeader(t, req, authToken, AUTHENTICATIONSCHEME, customer, -time.Minute)
			},
			responseCheck: func (t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func (t *testing.T) {
			server := newTestServer(t, nil)
			authUrl := "/authenticate"
			server.router.GET(
				authUrl,
				server.Authenticate,
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)
			responseRecorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authUrl, nil)
			require.NoError(t, err)
			tc.setUpAuth(t, req, server.token)
			server.router.ServeHTTP(responseRecorder, req)
			tc.responseCheck(t, responseRecorder)
		})
	}
}
