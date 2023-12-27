package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	// "strconv"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/token"
	"github.com/slamchillz/xchange/utils"
	"github.com/slamchillz/xchange/common"
	"github.com/go-playground/validator/v10"
)

type CoinSwapRequest struct {
	CoinName           string `json:"coin_name" binding:"required,coinname"`
	CoinAmountToSwap   float64 `json:"coin_amount_to_swap" binding:"required,numeric,gt=0"`
	Network            string `json:"network" binding:"required,oneof=BTC BSC TRON ETH,network"`
	PhoneNumber        string `json:"phone_number" binding:"required,phonenumber"`
	BankAccName        string `json:"bank_acc_name" binding:"required"`
	BankAccNumber      string `json:"bank_acc_number" binding:"required"`
	BankCode           string `json:"bank_code" binding:"required"`
}

type CoinSwapResponse struct {
	ID                 int32          `json:"id"`
	CoinName           string         `json:"coin_name"`
	CoinAmountToSwap   string         `json:"coin_amount_to_swap"`
	Network            string         `json:"network"`
	PhoneNumber        string         `json:"phone_number"`
	CoinAddress        string         `json:"coin_address"`
	TransactionRef     string         `json:"transaction_ref"`
	TransactionStatus  string         `json:"transaction_status"`
	CurrentUsdtNgnRate string         `json:"current_usdt_ngn_rate"`
	CustomerID         int32          `json:"customer_id"`
	NgnEquivalent      string         `json:"ngn_equivalent"`
	BankAccName        string         `json:"bank_acc_name"`
	BankAccNumber      string         `json:"bank_acc_number"`
	BankCode           string         `json:"bank_code"`
	CreatedAt		   time.Time      `json:"created_at"`
}

func (server *Server) CoinSwap(ctx *gin.Context) {
	var req CoinSwapRequest
	err := ctx.ShouldBindJSON(&req);
	// var bankErrChannel = make(chan error)
	var amountErrChannel = make(chan error)
	// go validateBankDetails(server.config, req.BankCode, req.BankAccNumber, bankErrChannel)
	// What if there is an error in the request fields coinname and coinamounttoswap
	// calling this function before validating the request fields could cause problems
	go validateCoinAmountToSwap(req.CoinName, req.CoinAmountToSwap, amountErrChannel)
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
	err = <-amountErrChannel
	if err != nil {
		reqErr["coin_amount_to_swap"] = err.Error()
		// close(amountErrChannel)
	}
	if len(reqErr) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": reqErr})
		return
	}
	customerPayload := ctx.MustGet(AUTHENTICATIONCONTEXTKEY).(*token.Payload) // revise this later
	arg := db.GetPendingNetworkTransactionParams{
		CustomerID: customerPayload.CustomerID,
		Network: req.Network,
		TransactionStatus: "PENDING",
	}
	checkPendingStart := time.Now()
	count, err := server.storage.GetPendingNetworkTransaction(context.Background(), arg)
	checkPendingEnd := time.Since(checkPendingStart)
	logger.Info().Msgf("check pending time: %v", checkPendingEnd)
	if err != nil {
		logger.Error().Err(err).Msg("error checking pending transaction")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if count > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("You have a pending %s transaction on %s network", req.CoinName, req.Network)})
		return
	}
	getAddressStart := time.Now()
	address, err := common.GetSwapAddress(server.config, server.storage, req.CoinName, req.Network, arg.CustomerID)
	getAddressEnd := time.Since(getAddressStart)
	logger.Info().Msgf("get address time: %v", getAddressEnd)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Error().Err(err).Msg("error getting address")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	currentUsdtRate := utils.RandomCoinswapRate()
	ngnEquivalent := req.CoinAmountToSwap * currentUsdtRate
	if strings.ToUpper(req.CoinName) == "BTC" {
		ngnEquivalent = ngnEquivalent * BTCUSDT
	}
	swapDetails, err := server.storage.CreateSwap(context.Background(), db.CreateSwapParams{
		CoinName: req.CoinName,
		CoinAmountToSwap: fmt.Sprintf("%f", req.CoinAmountToSwap),
		Network: req.Network,
		PhoneNumber: req.PhoneNumber,
		CoinAddress: address,
		TransactionRef: utils.RandomString(15),
		TransactionStatus: "PENDING",
		CurrentUsdtNgnRate: fmt.Sprintf("%f", currentUsdtRate),
		CustomerID: arg.CustomerID,
		NgnEquivalent: fmt.Sprintf("%f", ngnEquivalent),
		BankAccName: req.BankAccName,
		BankAccNumber: req.BankAccNumber,
		BankCode: req.BankCode,
	})
	if err != nil {
		logger.Error().Err(err).Msg("error creating swap")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	coinSwapResponse := CoinSwapResponse{
		ID: swapDetails.ID,
		CoinName: swapDetails.CoinName,
		CoinAmountToSwap: swapDetails.CoinAmountToSwap,
		Network: swapDetails.Network,
		PhoneNumber: swapDetails.PhoneNumber,
		CoinAddress: swapDetails.CoinAddress,
		TransactionRef: swapDetails.TransactionRef,
		TransactionStatus: swapDetails.TransactionStatus,
		CurrentUsdtNgnRate: swapDetails.CurrentUsdtNgnRate,
		CustomerID: swapDetails.CustomerID,
		NgnEquivalent: swapDetails.NgnEquivalent,
		BankAccName: swapDetails.BankAccName,
		BankAccNumber: swapDetails.BankAccNumber,
		BankCode: swapDetails.BankCode,
		CreatedAt: swapDetails.CreatedAt,
	}
	ctx.JSON(http.StatusOK, &coinSwapResponse)
}


type CoinSwapStatusUpdateRequest struct {
	TransactionStatus string `json:"transaction_status" binding:"required,oneof=CANCEL PAID"`
}


func (server *Server) CoinSwapStatusUpdate(ctx *gin.Context) {
	var err error
	var req CoinSwapStatusUpdateRequest
	var ve validator.ValidationErrors
	transactionRef, ok := ctx.Params.Get("ref")
	swap := make(chan db.Coinswap)
	swapErr := make(chan error)
	go func(store db.Store)  {
		doesSwapExist, err := store.GetOneCoinSwapTransaction(context.Background(), transactionRef)
		if err != nil {
			swapErr <- err
			swap <- db.Coinswap{}
			return
		}
		swapErr <- nil
		swap <- doesSwapExist
	}(server.storage)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing transaction reference"})
		return
	}
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		if errors.As(err, &ve) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unknown transaction status, expect CANCEL or PAID"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json request body"})
		}
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
	if err := <-swapErr; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		logger.Error().Err(err).Msg("error getting swap")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	swapDetails := <-swap
	if swapDetails.CustomerID != customerId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	// Consider race condition here. What if the transaction status is updated by background workers for the same transaction
	if swapDetails.CustomerAction == "CREATED" {
		arg := db.CoinSwapUpdateUserPaidParams{
			TransactionStatus: req.TransactionStatus,
			TransactionRef: transactionRef,
			CustomerID: customerId,
		}
		_, err = server.storage.CoinSwapUpdateUserPaid(context.Background(), arg)
		if err != nil {
			logger.Error().Err(err).Msg("error updating swap")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "transaction status updated successfully"})
		return
	}
	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "transaction is in a settled state"})
}
