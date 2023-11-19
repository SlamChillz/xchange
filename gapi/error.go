package gapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func authenticationError(err error) error {
	return status.Errorf(codes.Unauthenticated, "%s", err.Error())
}
