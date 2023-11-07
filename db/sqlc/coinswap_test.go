package db

import (
	"fmt"
	"context"
	"testing"

	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCoinSwap(t *testing.T) CreateSwapRow {
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
	require.Equal(t, arg.NgnEquivalent, coinswap.NgnEquivalent)
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
