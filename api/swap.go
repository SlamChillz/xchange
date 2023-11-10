package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	// "strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/token"
	"github.com/slamchillz/xchange/utils"
)

type CoinSwapRequest struct {
	CoinName           string `json:"coin_name" binding:"required"`
	CoinAmountToSwap   float64 `json:"coin_amount_to_swap" binding:"required,numeric,gt=0"`
	Network            string `json:"network" binding:"required,oneof=BTC BSC TRON ETH"`
	PhoneNumber        string `json:"phone_number" binding:"required,phonenumber"`
	BankAccName        string `json:"bank_acc_name" binding:"required"`
	BankAccNumber      string `json:"bank_acc_number" binding:"required"`
	BankCode           string `json:"bank_code" binding:"required"`
	FieldErrors		   map[string]string
	ServerErrored	   bool
	WaitGroup		   *sync.WaitGroup
	btcUsdRate		   float64
}

func (csr *CoinSwapRequest) validateCoinName()  {
	csr.CoinName = strings.ToUpper(csr.CoinName)
	if _, ok := BITPOWRCOINTICKER[csr.CoinName]; !ok {
		csr.FieldErrors["coin_name"] = "Invalid coin name"
	}
	csr.WaitGroup.Done()
}

func (csr *CoinSwapRequest) validateCoinAmountToSwap() {
	amount := csr.CoinAmountToSwap
	if csr.CoinName == "BTC" {
		amount = amount * csr.btcUsdRate
	}
	if amount < 20.00 {
		csr.FieldErrors["coin_amount_to_swap"] = "Coin amount to swap must be $20 or greater"
	}
	csr.WaitGroup.Done()
}

func (csr *CoinSwapRequest) validateNetwork() {
	csr.Network = strings.ToUpper(csr.Network)
	if network, ok := CHAINS[csr.CoinName]; !ok || network != csr.Network {
		csr.FieldErrors["network"] = "Invalid network selected or network not supported"
	}
	csr.WaitGroup.Done()
}

// func (csr *CoinSwapRequest) validatePhoneNumber() {
// 	phoneNumberLength := len(csr.PhoneNumber)
// 	if (strings.HasPrefix(csr.PhoneNumber, "+")  && phoneNumberLength != 14) || (strings.HasPrefix(csr.PhoneNumber, "0") && phoneNumberLength != 11)  {
// 		csr.FieldErrors["phone_number"] = "invalid phone number"
// 	}
// 	csr.WaitGroup.Done()
// }

func (csr *CoinSwapRequest) validateBankDetails() {
	config, err := utils.LoadConfig("../")
	if err != nil {
		csr.ServerErrored = true
	} else {
		rawData := map[string]string {
			"bank": csr.BankCode,
			"account": csr.BankAccNumber,
		}
		reqData := new(bytes.Buffer)
		err = json.NewEncoder(reqData).Encode(rawData)
		if err != nil {
			csr.ServerErrored = true
		} else {
			url, err := url.Parse(config.VALIDATE_BANK_URL)
			if err != nil {
				csr.ServerErrored = true
			} else {
				req := &http.Request{
					Method: http.MethodPost,
					URL: url,
					Header: map[string][]string {
						"Content-Type": {"application/json"},
						"X-API-KEY": {config.SHUTTER_PUBLIC_KEY},
					},
					Body: io.NopCloser(reqData),
				}
				response, err := http.DefaultClient.Do(req)
				if err != nil {
					csr.ServerErrored = true
				} else {
					defer response.Body.Close()
					if response.StatusCode != http.StatusOK {
						csr.FieldErrors["bank_acc_number"] = "Unable to verify account details"
					}
				}
			}
		}
	}
	csr.WaitGroup.Done()
}

func (csr *CoinSwapRequest) validate() {
	csr.WaitGroup = new(sync.WaitGroup)
	csr.WaitGroup.Add(4)
	go csr.validateBankDetails()
	go csr.validateCoinAmountToSwap()
	go csr.validateCoinName()
	go csr.validateNetwork()
	// go csr.validatePhoneNumber()
	csr.WaitGroup.Wait()
}

func (server *Server) CoinSwap(ctx *gin.Context) {
	var req CoinSwapRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.btcUsdRate = 34403.000000
	req.FieldErrors = make(map[string]string)
	req.WaitGroup = new(sync.WaitGroup)
	validationStart := time.Now()
	req.validate()
	validationEnd := time.Since(validationStart)
	logger.Println("validation time: ", validationEnd)
	if len(req.FieldErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": req.FieldErrors})
		return
	}
	customerPayload := ctx.MustGet(AUTHENTICATIONCONTEXTKEY).(*token.Payload)
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
	address, err := GetSwapAddress(server.config, server.storage, &req, arg.CustomerID)
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
		ngnEquivalent = ngnEquivalent * req.btcUsdRate
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
	ctx.JSON(http.StatusOK, &swapDetails)
}

func GenerateNewAddress(config utils.Config, customerId int32, label, asset, accountId string) (map[string]interface{}, error) {
	postData := map[string]string {
		"label": label,
		"asset": asset,
		"accountId": accountId,
	}
	reqData := new(bytes.Buffer)
	err := json.NewEncoder(reqData).Encode(postData)
	if err != nil {
		return nil, err
	}
	url, err := url.Parse("https://developers.bitpowr.com/api/v1/addresses")
	if err != nil {
		return nil, err
	}
	req := http.Request{
		Method: http.MethodPost,
		URL: url,
		Header: map[string][]string {
			"Content-Type": {"application/json"},
			"Authorization": {fmt.Sprintf("Bearer %s", config.BITPOWR_API_KEY)},
		},
		Body: io.NopCloser(reqData),
	}
	response, err := http.DefaultClient.Do(&req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to generate new %s address", asset)
	}
	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	data = data["data"].(map[string]interface{})
	return data, nil
}


func GetSwapAddress(config utils.Config, storage db.Store, req *CoinSwapRequest, customerId int32) (address string, err error) {
	var nullAddress sql.NullString
	var newAaddressData map[string]interface{}
	switch req.Network {
	case "BTC":
		nullAddress, err = storage.GetBtcAddress(context.Background(), customerId)
	case "BSC":
		nullAddress, err = storage.GetUsdtBscAddress(context.Background(), customerId)
	case "TRON":
		nullAddress, err = storage.GetUsdtTronAddress(context.Background(), customerId)
	case "ETH":
		nullAddress, err = storage.GetUsdtAddress(context.Background(), customerId)
	}
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return
		}
	}
	address = nullAddress.String
	if address == "" {
		label := req.CoinName + "_" + req.Network
		asset := req.CoinName
		accountId := config.BITPOWR_ACCOUNT_ID
		newAaddressData, err = GenerateNewAddress(config, customerId, label, asset, accountId)
		if err != nil {
			return
		}
	}
	address = newAaddressData["address"].(string)
	addressUid :=newAaddressData["uid"].(string)
	network := newAaddressData["chain"].(string)
	switch req.Network {
	case "BTC":
		_, err = storage.InsertNewBtcAddress(context.Background(), db.InsertNewBtcAddressParams{
			BtcAddress: sql.NullString{String: address, Valid: true},
			BtcAddressUid: sql.NullString{String: addressUid, Valid: true},
			BtcNetwork: sql.NullString{String: network, Valid: true},
			CustomerID: customerId,
		})
	case "BSC":
		_, err = storage.InsertNewUsdtBscAddress(context.Background(), db.InsertNewUsdtBscAddressParams{
			UsdtBscAddress: sql.NullString{String: address, Valid: true},
			UsdtBscAddressUid: sql.NullString{String: addressUid, Valid: true},
			UsdtBscNetwork: sql.NullString{String: network, Valid: true},
			CustomerID: customerId,
		})
	case "TRON":
		_, err = storage.InsertNewUsdtTronAddress(context.Background(), db.InsertNewUsdtTronAddressParams{
			UsdtTronAddress: sql.NullString{String: address, Valid: true},
			UsdtTronAddressUid: sql.NullString{String: addressUid, Valid: true},
			UsdtTronNetwork: sql.NullString{String: network, Valid: true},
			CustomerID: customerId,
		})
	case "ETH":
		_, err = storage.InsertNewUsdtAddress(context.Background(), db.InsertNewUsdtAddressParams{
			UsdtAddress: sql.NullString{String: address, Valid: true},
			UsdtAddressUid: sql.NullString{String: addressUid, Valid: true},
			UsdtNetwork: sql.NullString{String: network, Valid: true},
			CustomerID: customerId,
		})
	}
	if err != nil {
		return
	}
	return address, nil
}
