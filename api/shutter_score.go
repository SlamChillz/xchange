package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"encoding/json"
	"os"
)

var (
	ErrVerifyBankDetails = errors.New("unable to verify bank details")
)

type ShutterBankResponseData struct {
	BankName string `json:"bank_name"`
	BankCode string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
	AccountName string `json:"account_name"`
}

type ShutterBankResponse struct {
	Status bool `json:"status"`
	StatusCode int `json:"statusCode"`
	Message string `json:"message"`
	Data ShutterBankResponseData `json:"data"`
}

// Consider making a config package to hold all the config variables
func ValidateBankDetailsFromShutterScore(bankDetails BankDetailsRequest) (*ShutterBankResponse, error) {
	resolveUrl := "https://api.shutterscore.io/v1/merchant/public/misc/banks/resolve"
	rawData := map[string]string {
		"bank": bankDetails.BankCode,
		"account": bankDetails.BankAccNumber,
	}
	reqData := new(bytes.Buffer)
	err := json.NewEncoder(reqData).Encode(rawData)
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(resolveUrl)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL: url,
		Header: map[string][]string {
			"Content-Type": {"application/json"},
			// "X-API-KEY": {config.SHUTTER_PUBLIC_KEY},
			"X-API-KEY": {os.Getenv("SHUTTER_PUBLIC_KEY")},
		},
		Body: io.NopCloser(reqData),
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var shutterResponse ShutterBankResponse
	err = json.NewDecoder(response.Body).Decode(&shutterResponse)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, ErrVerifyBankDetails
	}
	return &shutterResponse, nil
}
