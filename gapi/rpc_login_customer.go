package gapi

import (
	"context"
	"errors"
	"fmt"
	"database/sql"

	"github.com/slamchillz/xchange/pb"
	"github.com/slamchillz/xchange/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginCustomer(ctx context.Context, req *pb.LoginCustomerRequest) (*pb.LoginCustomerResponse, error) {
	customer, err := server.storage.GetCustomerByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "customer not found")
		}
		// log the error
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal server error: %v", err.Error()))
	}
	err = utils.CheckPassword(customer.Password, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	accessToken, _, err := server.token.CreateToken(customer.ID, server.config.JWT_ACCESS_TOKEN_DURATION)
	if err != nil {
		// log the error
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal server error: %v", err.Error()))
	}
	resp := &pb.LoginCustomerResponse{
		AccessToken: accessToken,
		Customer: convertCustomerResponse(customer),
	}
	return resp, nil
}
