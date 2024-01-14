package types

import (
	"time"
)

type CreateCustomerRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3,max=50,alpha"`
	LastName string `json:"last_name" binding:"required,min=3,max=50,alpha"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	Verified bool
}

type CustomerResponse struct {
	ID 				int32 		   `json:"id"`
	FirstName 		string 		   `json:"first_name"`
	LastName 		string 		   `json:"last_name"`
	Email 			string 		   `json:"email"`
	IsActive        bool           `json:"is_active"`
	IsStaff         bool           `json:"is_staff"`
	IsSupercustomer bool           `json:"is_supercustomer"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type EmailSignupVerificationRequest struct {
	OTP string `json:"otp" binding:"required,number"`
}

type GoogleAuthRequest struct {
	Token string `json:"token" binding:"required"`
}

type GoogleUserInfo struct {
	Sub string `json:"sub"`
	Email string `json:"email"`
	EmailVerified bool `json:"email_verified"`
	FirstName string `json:"given_name"`
	LastName string `json:"family_name"`
	Picture string `json:"picture"`
	Locale string `json:"locale"`
}
