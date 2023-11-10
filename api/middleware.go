package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)


const (
	AUTHENTICATIONHEADER = "authorization"
	AUTHENTICATIONSCHEME = "bearer"
	AUTHENTICATIONCONTEXTKEY = "authenticated_customer"
)

func (server *Server) Authenticate(ctx *gin.Context) {
	authHeader := ctx.GetHeader(AUTHENTICATIONHEADER)
	if len(authHeader) <= len(AUTHENTICATIONSCHEME) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or incomplete authentication header"})
		return
	}
	authHeaderValues := strings.Fields(authHeader)
	if len(authHeaderValues) != 2 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unrecognized authentication header format"})
		return
	}
	if strings.ToLower(authHeaderValues[0]) != AUTHENTICATIONSCHEME {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unsupported authentication scheme"})
		return
	}
	accessToken := authHeaderValues[1]
	customer, err := server.token.VerifyToken(accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		return
	}
	ctx.Set(AUTHENTICATIONCONTEXTKEY, customer)
	ctx.Next()
}
