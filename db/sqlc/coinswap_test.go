package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCoinSwap(t *testing.T) Coinswap {
	customer := createRandomCustomer(t)
	network := utils.RandomCoinNetwork()
	coinaddress := utils.RandomCoinAddress(network).String
	coinswapamount := utils.RandomCoinSwapAmount()
	coinswaprate := utils.RandomCoinswapRate()
	arg := CreateSwapParams{
		CoinName: utils.RandomName(),
		CoinAmountToSwap: fmt.Sprintf("%.8f", coinswapamount),
		Network: network,
		PhoneNumber: utils.RandomPhoneNumber(),
		CoinAddress: coinaddress,
		TransactionRef: utils.RandomString(15),
		CurrentUsdtNgnRate: fmt.Sprintf("%v", coinswaprate),
		CustomerID: customer.ID,
		NgnEquivalent: fmt.Sprintf("%.8f", coinswapamount * coinswaprate),
		// PayoutStatus: utils.RandomPayoutStatus(),
		BankAccName: utils.RandomBankName(),
		BankAccNumber: utils.RandomBankAccount(),
		BankCode: utils.RandomBankCode(),
	}
	coinswap, err := testQueries.CreateSwap(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, coinswap)
	require.Equal(t, arg.CoinName, coinswap.CoinName)
	require.Equal(t, arg.CoinAmountToSwap, coinswap.CoinAmountToSwap)
	require.Equal(t, arg.Network, coinswap.Network)
	require.Equal(t, arg.PhoneNumber, coinswap.PhoneNumber)
	require.Equal(t, arg.CoinAddress, coinswap.CoinAddress)
	require.Equal(t, arg.TransactionRef, coinswap.TransactionRef)
	require.Equal(t, arg.CurrentUsdtNgnRate, coinswap.CurrentUsdtNgnRate)
	require.Equal(t, arg.CustomerID, coinswap.CustomerID)
	require.Equal(t, arg.BankAccName, coinswap.BankAccName)
	require.Equal(t, arg.BankAccNumber, coinswap.BankAccNumber)
	require.Equal(t, arg.BankCode, coinswap.BankCode)
	require.NotZero(t, coinswap.ID)
	require.NotZero(t, coinswap.CreatedAt)
	return coinswap
}

func TestCreateSwap(t *testing.T) {
	createRandomCoinSwap(t)
}

func updateRandomCoinSwapWithBitpowrInfo(t *testing.T) Coinswap {
	coinswap := createRandomCoinSwap(t)
	transactionAmountInDollars, err := strconv.ParseFloat(coinswap.CoinAmountToSwap, 64)
	require.NoError(t, err)
	if coinswap.CoinName == "BTC" {
		transactionAmountInDollars = transactionAmountInDollars * 32719.50
	}
	usdtNGNRate, _ := strconv.ParseFloat(coinswap.CurrentUsdtNgnRate, 64)
	transactionAmountInNgn := fmt.Sprintf("%f", transactionAmountInDollars * usdtNGNRate)
	arg := UpdateSwapWithBitpowrInfoParams {
		BitpowrRef: sql.NullString{
			String: utils.RandomString(15),
			Valid: true,
		},
		TransAddress: sql.NullString{
			String: coinswap.CoinAddress,
			Valid: true,
		},
		TransAmount: sql.NullString{
			String: coinswap.CoinAmountToSwap,
			Valid: true,
		},
		TransChain: sql.NullString{
			String: coinswap.Network,
			Valid: true,
		},
		TransHash: sql.NullString{
			String: utils.RandomString(15),
			Valid: true,
		},
		TransactionStatus: "SUCCESS",
		TransAmountNgn: sql.NullString{
			String: transactionAmountInNgn,
			Valid: true,
		},
		CoinName: coinswap.CoinName,
		CoinAddress: coinswap.CoinAddress,
	}
	updatedCoinswap, err := testQueries.UpdateSwapWithBitpowrInfo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedCoinswap)
	require.Equal(t, arg.BitpowrRef, updatedCoinswap.BitpowrRef)
	require.Equal(t, arg.TransAddress, updatedCoinswap.TransAddress)
	require.Equal(t, arg.TransAmount, updatedCoinswap.TransAmount)
	require.Equal(t, arg.TransChain, updatedCoinswap.TransChain)
	require.Equal(t, arg.TransHash, updatedCoinswap.TransHash)
	require.Equal(t, arg.TransactionStatus, updatedCoinswap.TransactionStatus)
	require.Equal(t, arg.TransAmountNgn, updatedCoinswap.TransAmountNgn)

	require.Equal(t, coinswap.CoinName, updatedCoinswap.CoinName)
	require.Equal(t, coinswap.CoinAddress, updatedCoinswap.CoinAddress)
	require.Equal(t, coinswap.CoinAmountToSwap, updatedCoinswap.CoinAmountToSwap)
	require.Equal(t, coinswap.Network, updatedCoinswap.Network)
	require.Equal(t, coinswap.PhoneNumber, updatedCoinswap.PhoneNumber)
	require.Equal(t, coinswap.TransactionRef, updatedCoinswap.TransactionRef)
	require.Equal(t, coinswap.CurrentUsdtNgnRate, updatedCoinswap.CurrentUsdtNgnRate)
	require.Equal(t, coinswap.CustomerID, updatedCoinswap.CustomerID)
	require.Equal(t, coinswap.BankAccName, updatedCoinswap.BankAccName)
	require.Equal(t, coinswap.BankAccNumber, updatedCoinswap.BankAccNumber)
	require.Equal(t, coinswap.BankCode, updatedCoinswap.BankCode)
	require.Equal(t, coinswap.PayoutStatus, updatedCoinswap.PayoutStatus)
	require.Equal(t, coinswap.ID, updatedCoinswap.ID)
	require.Equal(t, coinswap.CreatedAt, updatedCoinswap.CreatedAt)
	// require.NotEqual(t, coinswap.UpdatedAt, updatedCoinswap.UpdatedAt)
	return updatedCoinswap
}

func TestUpdateSwapWithBitpowrInfo(t *testing.T) {
	updateRandomCoinSwapWithBitpowrInfo(t)
}

func TestUpdateSwapWithShutterInfo(t *testing.T) {
	coinSwap := updateRandomCoinSwapWithBitpowrInfo(t)
	arg := UpdateSwapWithShutterInfoParams {
		PayoutStatus: sql.NullString{
			String: "SUCCESS",
			Valid: true,
		},
		AdminTransAmount: coinSwap.TransAmountNgn,
		AdminTransFee: sql.NullString{
			String: "0.00",
			Valid: true,
		},
		AdminTransRef: sql.NullString{
			String: utils.RandomString(15),
			Valid: true,
		},
		TransactionRef: coinSwap.TransactionRef,
	}
	updatedCoinswap, err := testQueries.UpdateSwapWithShutterInfo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedCoinswap)
	require.Equal(t, arg.PayoutStatus.String, updatedCoinswap.PayoutStatus)
	require.Equal(t, arg.AdminTransAmount, updatedCoinswap.AdminTransAmount)
	require.Equal(t, arg.AdminTransFee, updatedCoinswap.AdminTransFee)
	require.Equal(t, arg.AdminTransRef, updatedCoinswap.AdminTransRef)
	require.Equal(t, arg.TransactionRef, updatedCoinswap.TransactionRef)

	require.Equal(t, coinSwap.CoinName, updatedCoinswap.CoinName)
	require.Equal(t, coinSwap.CoinAddress, updatedCoinswap.CoinAddress)
	require.Equal(t, coinSwap.CoinAmountToSwap, updatedCoinswap.CoinAmountToSwap)
	require.Equal(t, coinSwap.Network, updatedCoinswap.Network)
	require.Equal(t, coinSwap.PhoneNumber, updatedCoinswap.PhoneNumber)
	require.Equal(t, coinSwap.TransactionRef, updatedCoinswap.TransactionRef)
	require.Equal(t, coinSwap.CurrentUsdtNgnRate, updatedCoinswap.CurrentUsdtNgnRate)
	require.Equal(t, coinSwap.CustomerID, updatedCoinswap.CustomerID)
	require.Equal(t, coinSwap.BankAccName, updatedCoinswap.BankAccName)
	require.Equal(t, coinSwap.BankAccNumber, updatedCoinswap.BankAccNumber)
	require.Equal(t, coinSwap.BankCode, updatedCoinswap.BankCode)
	require.Equal(t, coinSwap.ID, updatedCoinswap.ID)
	require.Equal(t, coinSwap.CreatedAt, updatedCoinswap.CreatedAt)
	// require.Equal(t, coinSwap.UpdatedAt, updatedCoinswap.UpdatedAt)
}
