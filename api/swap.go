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
	logger.Println("check pending time: ", checkPendingEnd)
	if err != nil {
		logger.Println(err)
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
	logger.Println("get address time: ", getAddressEnd)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logger.Println(err)
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
		logger.Println(err)
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
