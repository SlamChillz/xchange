// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CoinSwapUpdateUserPaid(ctx context.Context, arg CoinSwapUpdateUserPaidParams) (Coinswap, error)
	CreateAdmindata(ctx context.Context, arg CreateAdmindataParams) (Admindatum, error)
	CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error)
	CreateSwap(ctx context.Context, arg CreateSwapParams) (Coinswap, error)
	GetBtcAddress(ctx context.Context, customerID int32) (sql.NullString, error)
	GetCustomerBankDetails(ctx context.Context, customerID int32) (Bankdetail, error)
	GetCustomerByEmail(ctx context.Context, email string) (Customer, error)
	GetCustomerByGoogleId(ctx context.Context, googleID sql.NullString) (Customer, error)
	GetCustomerById(ctx context.Context, id int32) (Customer, error)
	GetCustomerByPhoneNumber(ctx context.Context, phone sql.NullString) (Customer, error)
	GetOneCoinSwapTransaction(ctx context.Context, transactionRef string) (Coinswap, error)
	GetPendingNetworkTransaction(ctx context.Context, arg GetPendingNetworkTransactionParams) (int64, error)
	GetUsdtAddress(ctx context.Context, customerID int32) (sql.NullString, error)
	GetUsdtBscAddress(ctx context.Context, customerID int32) (sql.NullString, error)
	GetUsdtTronAddress(ctx context.Context, customerID int32) (sql.NullString, error)
	InsertCustomerBankDetails(ctx context.Context, arg InsertCustomerBankDetailsParams) ([]Bankdetail, error)
	InsertNewBtcAddress(ctx context.Context, arg InsertNewBtcAddressParams) (Customerasset, error)
	InsertNewUsdtAddress(ctx context.Context, arg InsertNewUsdtAddressParams) (Customerasset, error)
	InsertNewUsdtBscAddress(ctx context.Context, arg InsertNewUsdtBscAddressParams) (Customerasset, error)
	InsertNewUsdtTronAddress(ctx context.Context, arg InsertNewUsdtTronAddressParams) (Customerasset, error)
	ListAllCoinSwapTransactions(ctx context.Context, arg ListAllCoinSwapTransactionsParams) ([]Coinswap, error)
	RecentUsdtNgnRate(ctx context.Context, usdtNgnRate string) (Usdtngnrate, error)
	ResetCustomerPassword(ctx context.Context, arg ResetCustomerPasswordParams) (Customer, error)
	UpdateBtcAddress(ctx context.Context, arg UpdateBtcAddressParams) (Customerasset, error)
	UpdateCustomerPassword(ctx context.Context, arg UpdateCustomerPasswordParams) (Customer, error)
	UpdateSwapWithBitpowrInfo(ctx context.Context, arg UpdateSwapWithBitpowrInfoParams) (Coinswap, error)
	UpdateSwapWithShutterInfo(ctx context.Context, arg UpdateSwapWithShutterInfoParams) (Coinswap, error)
	UpdateUsdtAddress(ctx context.Context, arg UpdateUsdtAddressParams) (Customerasset, error)
	UpdateUsdtBscAddress(ctx context.Context, arg UpdateUsdtBscAddressParams) (Customerasset, error)
	UpdateUsdtTronAddress(ctx context.Context, arg UpdateUsdtTronAddressParams) (Customerasset, error)
}

var _ Querier = (*Queries)(nil)
