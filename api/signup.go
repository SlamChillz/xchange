package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/utils"
)

type CreateCustomerRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3,max=50,alpha"`
	LastName string `json:"last_name" binding:"required,min=3,max=50,alpha"`
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
	var err error
	var req CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.PhoneNumber[0] == '0' {
		req.PhoneNumber = "+234" + req.PhoneNumber[1:]
	}
	// This is an expensive operation, should be done in a goroutine/queue
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// TODO: Validate phone number and email to avoid pk increment on db error
	customer, err := server.storage.GetCustomerByEmail(ctx, req.Email)
	if err != nil {
		// log the error
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error().Err(err).Msg("error getting customer by email")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		customer, err = server.storage.GetCustomerByPhoneNumber(ctx, sql.NullString{
			String: req.PhoneNumber,
			Valid: true,
		})
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logger.Error().Err(err).Msg("error getting customer by phone number")
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		} else if customer != (db.Customer{}) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "phone number has been taken"})
			return
		}
	} else if customer != (db.Customer{}) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email has been taken"})
		return
	}
	// Create the customer
	// Save to redis and send it to a queue
	customer, err = server.storage.CreateCustomer(ctx, db.CreateCustomerParams{
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Phone: sql.NullString{
			String: req.PhoneNumber,
			Valid: true,
		},
		Password: sql.NullString{
			String: hashedPassword,
			Valid: true,
		},
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505"{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email or phone number already exists"})
		} else {
			// log the error
			logger.Error().Err(err).Msg("error creating customer")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	otp := utils.OTP.GenerateOTP()
	_, err = server.redisClient.Set(ctx, otp, customer.Email, 5 * time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
		PhoneNumber: customer.Phone.String,
		IsActive: customer.IsActive,
		IsStaff: customer.IsStaff,
		IsSupercustomer: customer.IsSupercustomer,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}


type GoogleAuthRequest struct {
	Token string `json:"token" binding:"required"`
}


func (server *Server) GoogleSignUp(ctx *gin.Context) {
	var req GoogleAuthRequest
	var ve validator.ValidationErrors
	err := ctx.ShouldBindJSON(&req)
	reqErr := CreateSwapError{}
	if err != nil {
		if errors.As(err, &ve) {
			for _, e := range ve {
				field := e.Field()
				key, value := genrateFieldErrorMessage(field)
				if key != "" {
					reqErr[key] = value
				}
			}
		} else {
			logger.Error().Err(err).Msg("error binding request body")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json request body"})
			return
		}
	}
	if len(reqErr) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": reqErr})
		return
	}
	userInfo, err := getUserInfo(req.Token)
	if err != nil {
		logger.Error().Err(err).Msg("error getting token info")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}
	if !userInfo.EmailVerified {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email not verified"})
		return
	}
	_, err = server.storage.CreateCustomer(ctx, db.CreateCustomerParams{
		FirstName: userInfo.FirstName,
		LastName: userInfo.LastName,
		Email: userInfo.Email,
		Photo: sql.NullString{
			String: userInfo.Picture,
			Valid: true,
		},
		GoogleID: sql.NullString{
			String: userInfo.Sub,
			Valid: true,
		},
		Password: sql.NullString{},
		Phone: sql.NullString{},
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			logger.Error().Interface("pgErr", pgErr).Err(err).Msg("error creating customer via google")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "email has been taken"})
		} else {
			// log the error
			logger.Error().Err(err).Msg("error creating customer via google")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
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

func getUserInfo(acessToken string) (GoogleUserInfo, error) {
	var userInfo GoogleUserInfo
	url := "https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + acessToken
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return userInfo, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return userInfo, err
	}
	defer resp.Body.Close()
	byteResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return userInfo, err
	}
	err = json.Unmarshal(byteResponse, &userInfo)
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}
