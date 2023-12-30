package api

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	// "github.com/slamchillz/xchange/utils"
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
	case "OldPassword":
		return "old_password", "incorrect old password"
	case "NewPassword":
		return "new_password", "new password must differ from old password and must be at least 8 characters"
	case "ConfirmNewPassword":
		return "confirm_new_password", "new password does not match confirm new password"
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

var validateOldPassword validator.Func = func(fieldLevel validator.FieldLevel) bool {
	parent := fieldLevel.Parent().Interface().(ChangePasswordRequest)
	logger.Info().Interface("payload", parent).Msg("validating old password")
	if oldPassword, ok := fieldLevel.Field().Interface().(string); ok {
		return len(oldPassword) >= 8
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
	} else {
		amountErrChannel <- nil
	}
	close(amountErrChannel)
}

// func validateBankDetails(config utils.Config, bankCode string, bankAccNumber string, bankErrChannel chan<- error) {
// 	bankDetails := BankDetailsRequest{
// 		BankCode: bankCode,
// 		BankAccNumber: bankAccNumber,
// 	}
// 	_, err := ValidateBankDetailsFromShutterScore(bankDetails)
// 	if err != nil {
// 		logger.Error().Err(err).Msg("error validating bank details")
// 		bankErrChannel <- err
// 	}
// 	bankErrChannel <- nil
// }
