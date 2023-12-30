package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/token"
	"github.com/slamchillz/xchange/utils"
)


type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,gte=8,oldpassword"`
	NewPassword string `json:"new_password" binding:"required,gte=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,eqfield=NewPassword"`
}


func (server *Server) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	var ve validator.ValidationErrors
	reqErr := CreateSwapError{}
	err := ctx.ShouldBindJSON(&req)
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
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid json request body"})
			return
		}
	}
	if len(reqErr) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": reqErr})
		return
	}
	payload, ok := ctx.Get(AUTHENTICATIONCONTEXTKEY)
	if !ok{
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
	customer, err := server.storage.GetCustomerById(ctx, customerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		// log the error
		logger.Error().Err(err).Int32("customerid", customerId).Msg("Could not find authenticated customer in database")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = utils.CheckPassword(customer.Password, req.OldPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect old password"})
		return
	}
	newPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logger.Error().Err(err).Int32("customerid", customerId).Str("new password", req.NewPassword).Msg("Could not hash new password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	_, err = server.storage.UpdateCustomerPassword(ctx, db.UpdateCustomerPasswordParams{
		ID: customerId,
		Password: newPassword,
	})
	if err != nil {
		logger.Error().Err(err).Int32("customerid", customerId).Str("new password", req.NewPassword).Msg("Could not update new password into database")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	deleted, err := server.redisClient.Delete(ctx, customerPayload.ID)
	if err != nil || deleted == 0 {
		logger.Error().Err(err).Str("redis key", customerPayload.ID).Msg("error deleting user token info from redis store after password change")
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"status": "SUCCESS",
		"message": "login with your new password",
	})
}
