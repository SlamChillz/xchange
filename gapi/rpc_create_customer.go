package gapi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/pb"
	"github.com/slamchillz/xchange/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	phoneNumber := req.GetPhoneNumber()
	if phoneNumber[0] == '0' {
		phoneNumber = "+234" + phoneNumber[1:]
	}
	// TODO: Validate phone number and email to avoid pk increment on db error
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
	customer, err := server.storage.CreateCustomer(ctx, db.CreateCustomerParams{
		FirstName: req.GetFirstName(),
		LastName: req.GetLastName(),
		Email: req.GetEmail(),
		Phone: sql.NullString{
			String: phoneNumber,
			Valid: true,
		},
		Password: sql.NullString{
			String: hashedPassword,
			Valid: true,
		},
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505"{
			return nil, status.Errorf(codes.AlreadyExists, "user with the email or phone number already exists")
		}
		// log the error
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal server error: %v", err.Error()))
	}
	resp := &pb.CreateCustomerResponse{
		Customer: convertCustomerResponse(customer),
	}
	// TODO: Send verification email to customer
	return resp, nil
}
