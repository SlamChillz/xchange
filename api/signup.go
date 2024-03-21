package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/notification/mail"
	apiTypes "github.com/slamchillz/xchange/types/api"
	"github.com/slamchillz/xchange/utils"
	"github.com/slamchillz/xchange/worker"
)

func (server *Server) CreateCustomer(ctx *gin.Context) {
	defer func() {
		panicValue := recover()
		if panicValue != nil {
			logger.Error().Interface("panic", panicValue).Msg("panic occurred")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()
	var err error
	var req apiTypes.CreateCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// This is an expensive operation, should be done in a goroutine/queue
	// hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	req.Verified = false
	// TODO: Validate phone number and email to avoid pk increment on db error
	customer, err := server.storage.GetCustomerByEmail(ctx, req.Email)
	if err != nil {
		// log the error
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error().Err(err).Msg("error getting customer by email")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	} else if customer != (db.Customer{}) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email has been taken"})
		return
	}
	// otp := utils.OTP.GenerateOTP()
	jsonReqData, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).Msg("failed to marshal request body")
		return
	}
	_, err = server.redisClient.Set(ctx, "signup:"+req.Email, jsonReqData, 0 * time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).Msg("failed to save new customer info in redis db store")
		return
	}
	// Create the customer
	err = server.taskDistributor.DistributeTaskStoreNewCustomer(
		context.Background(),
		worker.PayloadStoreNewCustomer{
			Data: jsonReqData,
			TaskOptions: nil,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).Msg("failed to distribute store-new-customer-in-db task")
		return
	}
	otp := utils.OTP.GenerateOTP()
	_, err = server.redisClient.Set(ctx, "signup:"+otp, req.Email, 5 * time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).Msg("failed to set otp in redis db store")
		return
	}
	err = server.taskDistributor.DistributeTaskSendMail(
		context.Background(),
		worker.PayloadSendMail{
			MailReciever: req.Email,
			Template: mail.TemplateVerificationEmail,
			TemplateData: mail.TemplateDataVerificationEmail{
				Email: req.Email,
				FirstName: req.FirstName,
				Otp: otp,
			},
			MailType: mail.MailTypeVerificationEmail,
			TaskOptions: nil,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).Msg("failed to distribute send-mail task for new customer verification")
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{})
}

func CreateCustomerResponse(customer db.Customer) apiTypes.CustomerResponse {
	return apiTypes.CustomerResponse{
		ID: customer.ID,
		FirstName: customer.FirstName,
		LastName: customer.LastName,
		Email: customer.Email,
		IsActive: customer.IsActive,
		IsStaff: customer.IsStaff,
		IsSupercustomer: customer.IsSupercustomer,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}


func (server *Server) EmailSignupVerification(ctx *gin.Context) {
	var email string
	var req apiTypes.EmailSignupVerificationRequest
	err := ctx.ShouldBindJSON(&req);
	reqErr := CreateSwapError{}
	var ve validator.ValidationErrors
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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json request body"})
			return
		}
	}
	if len(reqErr) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": reqErr})
		return
	}
	err = server.redisClient.ScanDel(ctx, "signup:"+req.OTP, &email)
	if err != nil {
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "otp has expired or is invalid"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			logger.Error().Err(err).
				Str("otp", req.OTP).
				Msg("failed to retrieve otp from redis db store")
		}
		return
	}
	var customerInfo apiTypes.CreateCustomerRequest
	redisCustomerInfo, _ := server.redisClient.Get(ctx, "signup:"+email)
	err = json.Unmarshal([]byte(redisCustomerInfo), &customerInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).
			Str("email", email).
			Msg("failed to unmarshal customer info from redis db store for signup verification")
		return
	}
	fmt.Printf("%+v\n", customerInfo)
	customerInfo.Verified = true
	updatedCustomerInfo, err := json.Marshal(customerInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).
			Str("email", email).
			Msg("failed to marshal verified customer info during signup verification")
		return
	}
	_, err = server.redisClient.Set(ctx, "signup:"+email, updatedCustomerInfo, 0 * time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).
			Str("email", email).
			Msg("failed to set verified customer info in redis db store during signup verification")
		return
	}
	taskContext := context.Background()
	err = server.taskDistributor.DistributeTaskStoreNewCustomer(
		taskContext,
		worker.PayloadStoreNewCustomer{
			Data: updatedCustomerInfo,
			TaskOptions: nil,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logger.Error().Err(err).
			Str("email", email).
			Msg("failed to distribute verify-new-customer-in-db task during signup verification")
		return
	}
	err = server.taskDistributor.DistributeTaskSendMail(
		taskContext,
		worker.PayloadSendMail{
			MailReciever: customerInfo.Email,
			Template: mail.TemplateWelcomeEmail,
			TemplateData: mail.TemplateDataWelcomeEmail{
				Email: customerInfo.Email,
				FirstName: customerInfo.FirstName,
			},
			MailType: mail.MailTypeWelcomeEmail,
			TaskOptions: nil,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "try again",
			"error": "internal server error",
		})
		logger.Error().Err(err).
			Str("customer-email", customerInfo.Email).
			Msg("failed to send welcome email after signup verification")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}


func (server *Server) GoogleSignUp(ctx *gin.Context) {
	var req apiTypes.GoogleAuthRequest
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


func getUserInfo(acessToken string) (apiTypes.GoogleUserInfo, error) {
	var userInfo apiTypes.GoogleUserInfo
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
