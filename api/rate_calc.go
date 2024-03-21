package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/slamchillz/xchange/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CoinNGNEquivalentRequest struct {
	CoinName string `json:"coin_name" binding:"required,coinname"`
	CoinAmountToSwap float64 `json:"coin_amount" binding:"required,numeric,gt=0"`
}


func (server Server) GetCoinNGNEquivalent(ctx *gin.Context) {
	var req CoinNGNEquivalentRequest
	var amountErrChannel = make(chan error)
	var ve validator.ValidationErrors
	err := ctx.ShouldBindJSON(&req)
	go validateCoinAmountToSwap(
		req.CoinName,
		req.CoinAmountToSwap,
		amountErrChannel,
	)
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
	if err = <-amountErrChannel; err != nil {
		reqErr["coin_amount"] = "Coin amount must be $20 or greater"
		ctx.JSON(http.StatusBadRequest, gin.H{"error": reqErr})
		return
	}
	// to be replaced with real time USDTNGN rate
	currentUsdtRate := utils.RandomCoinswapRate()
	ngnEquivalent := req.CoinAmountToSwap * currentUsdtRate
	if strings.ToUpper(req.CoinName) == "BTC" {
		ngnEquivalent = ngnEquivalent * BTCUSDT
	}
	ctx.JSON(http.StatusOK, gin.H{"total_coin_price_ngn": ngnEquivalent})
}
