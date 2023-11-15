package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slamchillz/xchange/utils"
)

type LoginCustomerRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8"`
}

type LoginCustomerResponse struct {
	AccessToken string `json:"access_token"`
	Customer CustomerResponse `json:"user"`
}

func (server *Server) LoginCustomer(ctx *gin.Context) {
	var req LoginCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer, err := server.storage.GetCustomerByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		// log the error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = utils.CheckPassword(customer.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	accessToken, _, err := server.token.CreateToken(customer.ID, server.config.JWT_ACCESS_TOKEN_DURATION)
	if err != nil {
		// log the error
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := LoginCustomerResponse{
		AccessToken: accessToken,
		Customer: CreateCustomerResponse(customer),
	}
	ctx.JSON(http.StatusOK, resp)
}
