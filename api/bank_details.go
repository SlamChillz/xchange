package api

import (
	"errors"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/token"
	"github.com/go-playground/validator/v10"
)

type BankDetailsRequest struct {
	BankAccNumber string `json:"bank_acc_number" binding:"required,len=10"`
	BankCode string `json:"bank_code" binding:"required,len=3"`
}


func (server Server) GetBankDetails(ctx *gin.Context) {
	defer func () {
		ctx.Request.Body.Close()
		panicValue := recover()
		if panicValue != nil {
			logger.Error().Msgf("panic: %v", panicValue)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()
	payload, ok := ctx.Get(AUTHENTICATIONCONTEXTKEY)
	if !ok {
		logger.Error().Msg("error getting customer payload from authentication context key")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	customerPayload, ok := payload.(*token.Payload)
	if !ok {
		logger.Error().Interface("payload", payload).Msg("error casting customer payload to token payload")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	customerId := customerPayload.CustomerID
	bankDetails, err := server.storage.GetCustomerBankDetails(ctx, customerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			emptyBankDetails := []db.Bankdetail{}
			ctx.JSON(http.StatusOK, emptyBankDetails)
			return
		}
		logger.Error().Err(err).Msg("error getting bank details")
		panic(err)
	}
	ctx.JSON(http.StatusOK, bankDetails)
}

func (server Server) AddBankDetails(ctx *gin.Context) {
	defer func () {
		ctx.Request.Body.Close()
		panicValue := recover()
		if panicValue != nil {
			logger.Error().Msgf("panic: %v", panicValue)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()
	var req BankDetailsRequest
	err := ctx.ShouldBindJSON(&req)
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
	customerPayload, ok := ctx.MustGet(AUTHENTICATIONCONTEXTKEY).(*token.Payload)
	if !ok {
		logger.Error().Msg("error casting customer payload to token payload")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	customerId := customerPayload.CustomerID
	// This surely needs to be refactored
	// This should be done asynchronously and pushed to the client through a websocket or push notification
	validBankDetails, err := ValidateBankDetailsFromShutterScore(req)
	if err != nil {
		if errors.Is(err, ErrVerifyBankDetails) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank details"})
			return
		}
		logger.Error().Err(err).Msg("error validating bank details")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	bankData := db.InsertCustomerBankDetailsParams{
		CustomerID: customerId,
		BankName: validBankDetails.Data.BankName,
		BankCode: validBankDetails.Data.BankCode,
		AccountNumber: validBankDetails.Data.AccountNumber,
		AccountName: validBankDetails.Data.AccountName,
	}
	bankDetails, err := server.storage.InsertCustomerBankDetails(ctx, bankData)
	if err != nil {
		logger.Error().Err(err).Msg("error creating bank details")
		panic(err)
	}
	ctx.JSON(http.StatusOK, bankDetails)
}
