package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/slamchillz/xchange/token"
	"google.golang.org/grpc/metadata"
)

const (
	AUTHENTICATIONHEADER = "authorization"
	AUTHENTICATIONSCHEME = "bearer"
)

func (server *Server) AuthenticateUser(ctx context.Context) (*token.Payload, error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("request metadata not found")
	}
	authValueList := metaData.Get(AUTHENTICATIONHEADER)
	bearerHeader := authValueList[0]
	if len(authValueList) == 0 || len(bearerHeader) <= len(AUTHENTICATIONSCHEME) {
		return nil, fmt.Errorf("missing or incomplete authentication header")
	}
	bearerHeaderList := strings.Fields(bearerHeader)
	if len(bearerHeaderList) != 2 {
		return nil, fmt.Errorf("unrecognized authentication header format")
	}
	scheme := bearerHeaderList[0]
	if strings.ToLower(scheme) != AUTHENTICATIONSCHEME {
		return nil, fmt.Errorf("unsupported authentication scheme")
	}
	accessToken := bearerHeaderList[1]
	payload, err := server.token.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}
	return payload, nil
}
