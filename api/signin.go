package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	err = utils.CheckPassword(customer.Password.String, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	accessToken, tokenPayload, err := server.token.CreateToken(customer.ID, server.config.JWT_ACCESS_TOKEN_DURATION)
	if err != nil {
		// log the error
		logger.Error().Err(err).Msg("error creating access token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	_, err = server.redisClient.Set(ctx, tokenPayload.ID, tokenPayload, server.config.JWT_ACCESS_TOKEN_DURATION)
	if err != nil {
		logger.Error().Interface("token_payload", tokenPayload).Err(err).Msg("unable to persist payload in redis store")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	resp := LoginCustomerResponse{
		AccessToken: accessToken,
		Customer: CreateCustomerResponse(customer),
	}
	ctx.JSON(http.StatusOK, resp)
}


func (server *Server) GoogleSignIn(ctx *gin.Context) {
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
	customer, err := server.storage.GetCustomerByGoogleId(ctx, sql.NullString{
		String: userInfo.Sub,
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		// log the error
		logger.Error().Err(err).Msg("Signin with Google: error getting customer by google id")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	accessToken, tokenPayload, err := server.token.CreateToken(customer.ID, server.config.JWT_ACCESS_TOKEN_DURATION)
	if err != nil {
		// log the error
		logger.Error().Err(err).Msg("error creating access token")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	_, err = server.redisClient.Set(ctx, tokenPayload.ID, tokenPayload, server.config.JWT_ACCESS_TOKEN_DURATION)
	if err != nil {
		logger.Error().Interface("token_payload", tokenPayload).Err(err).Msg("unable to persist payload in redis store")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	resp := LoginCustomerResponse{
		AccessToken: accessToken,
		Customer: CreateCustomerResponse(customer),
	}
	ctx.JSON(http.StatusOK, resp)
}
