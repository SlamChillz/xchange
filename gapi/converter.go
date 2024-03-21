package gapi

import (
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertCustomerResponse(customer db.Customer) *pb.Customer {
	return &pb.Customer{
		Id: int32(customer.ID),
		FirstName: customer.FirstName,
		LastName: customer.LastName,
		Email: customer.Email,
		Phone: customer.Phone.String,
		IsActive: customer.IsActive,
		IsStaff: customer.IsStaff,
		IsSupercustomer: customer.IsSupercustomer,
		CreatedAt: timestamppb.New(customer.CreatedAt),
		UpdatedAt: timestamppb.New(customer.UpdatedAt),
	}
}
