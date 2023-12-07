package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/slamchillz/xchange/utils"
)

type CreateSwapError map[string]string

var validatePhoneNumber validator.Func = func(fieldLevel validator.FieldLevel) bool {
	isPhoneNumber := regexp.MustCompile(`^(\+234|0)\d{10}$`)
	if phoneNumber, ok := fieldLevel.Field().Interface().(string); ok {
		return isPhoneNumber.MatchString(phoneNumber)
	}
	return false
}

func genrateFieldErrorMessage(field string) (string, string) {
	switch field {
	case "CoinName":
		return "coin_name", "Invalid coin name"
	case "Network":
		return "network", "Invalid network selected or network not supported"
	case "CoinAmountToSwap":
		return "coin_amount_to_swap", "Coin amount to swap must be $20 or greater"
	case "PhoneNumber":
		return "phone_number", "Invalid phone number"
	case "BankAccName":
		return "bank_acc_name", "Invalid bank account name"
	case "BankAccNumber":
		return "bank_acc_number", "Invalid bank account number"
	case "BankCode":
		return "bank_code", "Invalid bank code"
	}
	return "", ""
}

var validateCoinName validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if coinName, ok := fieldLevel.Field().Interface().(string); ok {
		fmt.Println(coinName)
		if _, ok := BITPOWRCOINTICKER[coinName]; ok {
			return true
		}
	}
	return false
}

var validateNetwork validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if network, ok := fieldLevel.Field().Interface().(string); ok {
		network = strings.ToUpper(network)
		if _, ok := NETWORKS[network]; ok {
			return true
		}
	}
	return false
}

func validateCoinAmountToSwap(coinName string, amountToSwap float64, amountErrChannel chan<- error) {
	amount := amountToSwap
	if coinName == "BTC" {
		amount = amount * BTCUSDT
	}
	if amount < 20.00 {
		amountErrChannel <- errors.New("coin amount to swap must be $20 or greater")
	}
	amountErrChannel <- nil
	close(amountErrChannel)
}

func validateBankDetails(config utils.Config, bankCode string, bankAccNumber string, bankErrChannel chan<- error) {
	rawData := map[string]string {
		"bank": bankCode,
		"account": bankAccNumber,
	}
	reqData := new(bytes.Buffer)
	err := json.NewEncoder(reqData).Encode(rawData)
	if err != nil {
		bankErrChannel <- err
	}
	url, err := url.Parse(config.VALIDATE_BANK_URL)
	if err != nil {
		bankErrChannel <- err
	}
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
		bankErrChannel <- err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusBadRequest {
		bankErrChannel <- errors.New("unable to verify account number")
	}
	bankErrChannel <- nil
}
