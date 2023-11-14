package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/utils"
)

type CreateCustomerRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3,max=50"`
	LastName string `json:"last_name" binding:"required,min=3,max=50"`
	Email string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,phonenumber"`
	Password string `json:"password" binding:"required,gte=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type CustomerResponse struct {
	ID 				int32 		   `json:"id"`
	FirstName 		string 		   `json:"first_name"`
	LastName 		string 		   `json:"last_name"`
	Email 			string 		   `json:"email"`
	PhoneNumber 	string 		   `json:"phone_number"`
	IsActive        bool           `json:"is_active"`
	IsStaff         bool           `json:"is_staff"`
	IsSupercustomer bool           `json:"is_supercustomer"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (server *Server) CreateCustomer(ctx *gin.Context) {
	var req CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.PhoneNumber[0] == '0' {
		req.PhoneNumber = "+234" + req.PhoneNumber[1:]
	}
	// TODO: Validate phone number and email to avoid pk increment on db error
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	customer, err := server.storage.CreateCustomer(ctx, db.CreateCustomerParams{
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Phone: req.PhoneNumber,
		Password: hashedPassword,
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505"{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or phone number already exists"})
		} else {
			// log the error
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	resp := CreateCustomerResponse(customer)
	// TODO: Send verification email to customer
	ctx.JSON(http.StatusOK, resp)
}

func CreateCustomerResponse(customer db.Customer) CustomerResponse {
	return CustomerResponse{
		ID: customer.ID,
		FirstName: customer.FirstName,
		LastName: customer.LastName,
		Email: customer.Email,
		PhoneNumber: customer.Phone,
		IsActive: customer.IsActive,
		IsStaff: customer.IsStaff,
		IsSupercustomer: customer.IsSupercustomer,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}
