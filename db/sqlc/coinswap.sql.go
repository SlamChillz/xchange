// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: coinswap.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const createSwap = `-- name: CreateSwap :one
INSERT INTO coinswap (
    coin_name, coin_amount_to_swap, network, phone_number,
    coin_address, transaction_ref, transaction_status, current_usdt_ngn_rate,
    customer_id, ngn_equivalent, bank_acc_name, bank_acc_number, bank_code
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12, $13
) RETURNING id, coin_name, coin_amount_to_swap, network, phone_number, coin_address, transaction_ref, transaction_status, current_usdt_ngn_rate, customer_id, ngn_equivalent, payout_status, bank_acc_name, bank_acc_number, bitpowr_ref, trans_address, trans_amount, trans_chain, trans_hash, bank_code, admin_trans_amount, admin_trans_fee, admin_trans_ref, admin_trans_uid, trans_amount_ngn, created_at, updated_at
`

type CreateSwapParams struct {
	CoinName           string `json:"coin_name"`
	CoinAmountToSwap   string `json:"coin_amount_to_swap"`
	Network            string `json:"network"`
	PhoneNumber        string `json:"phone_number"`
	CoinAddress        string `json:"coin_address"`
	TransactionRef     string `json:"transaction_ref"`
	TransactionStatus  string `json:"transaction_status"`
	CurrentUsdtNgnRate string `json:"current_usdt_ngn_rate"`
	CustomerID         int32  `json:"customer_id"`
	NgnEquivalent      string `json:"ngn_equivalent"`
	BankAccName        string `json:"bank_acc_name"`
	BankAccNumber      string `json:"bank_acc_number"`
	BankCode           string `json:"bank_code"`
}

func (q *Queries) CreateSwap(ctx context.Context, arg CreateSwapParams) (Coinswap, error) {
	row := q.db.QueryRowContext(ctx, createSwap,
		arg.CoinName,
		arg.CoinAmountToSwap,
		arg.Network,
		arg.PhoneNumber,
		arg.CoinAddress,
		arg.TransactionRef,
		arg.TransactionStatus,
		arg.CurrentUsdtNgnRate,
		arg.CustomerID,
		arg.NgnEquivalent,
		arg.BankAccName,
		arg.BankAccNumber,
		arg.BankCode,
	)
	var i Coinswap
	err := row.Scan(
		&i.ID,
		&i.CoinName,
		&i.CoinAmountToSwap,
		&i.Network,
		&i.PhoneNumber,
		&i.CoinAddress,
		&i.TransactionRef,
		&i.TransactionStatus,
		&i.CurrentUsdtNgnRate,
		&i.CustomerID,
		&i.NgnEquivalent,
		&i.PayoutStatus,
		&i.BankAccName,
		&i.BankAccNumber,
		&i.BitpowrRef,
		&i.TransAddress,
		&i.TransAmount,
		&i.TransChain,
		&i.TransHash,
		&i.BankCode,
		&i.AdminTransAmount,
		&i.AdminTransFee,
		&i.AdminTransRef,
		&i.AdminTransUid,
		&i.TransAmountNgn,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPendingNetworkTransaction = `-- name: GetPendingNetworkTransaction :one
SELECT COUNT(*) FROM coinswap WHERE customer_id = $1 AND network = $2 AND transaction_status = $3
`

type GetPendingNetworkTransactionParams struct {
	CustomerID        int32  `json:"customer_id"`
	Network           string `json:"network"`
	TransactionStatus string `json:"transaction_status"`
}

func (q *Queries) GetPendingNetworkTransaction(ctx context.Context, arg GetPendingNetworkTransactionParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPendingNetworkTransaction, arg.CustomerID, arg.Network, arg.TransactionStatus)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listAllCoinSwapTransactions = `-- name: ListAllCoinSwapTransactions :many
SELECT id, coin_name, coin_amount_to_swap, network, phone_number, coin_address, transaction_ref, transaction_status, current_usdt_ngn_rate, customer_id, ngn_equivalent, payout_status, bank_acc_name, bank_acc_number, bitpowr_ref, trans_address, trans_amount, trans_chain, trans_hash, bank_code, admin_trans_amount, admin_trans_fee, admin_trans_ref, admin_trans_uid, trans_amount_ngn, created_at, updated_at FROM coinswap
WHERE coin_name = ANY($1::varchar[])
AND transaction_status = ANY($2::varchar[])
AND network = ANY($3::varchar[])
AND created_at BETWEEN $4 AND $5
ORDER BY created_at DESC
LIMIT $7
OFFSET $6
`

type ListAllCoinSwapTransactionsParams struct {
	CoinName          []string  `json:"coin_name"`
	TransactionStatus []string  `json:"transaction_status"`
	Network           []string  `json:"network"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	Offset            int32     `json:"offset"`
	Limit             int32     `json:"limit"`
}

func (q *Queries) ListAllCoinSwapTransactions(ctx context.Context, arg ListAllCoinSwapTransactionsParams) ([]Coinswap, error) {
	rows, err := q.db.QueryContext(ctx, listAllCoinSwapTransactions,
		pq.Array(arg.CoinName),
		pq.Array(arg.TransactionStatus),
		pq.Array(arg.Network),
		arg.StartDate,
		arg.EndDate,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Coinswap{}
	for rows.Next() {
		var i Coinswap
		if err := rows.Scan(
			&i.ID,
			&i.CoinName,
			&i.CoinAmountToSwap,
			&i.Network,
			&i.PhoneNumber,
			&i.CoinAddress,
			&i.TransactionRef,
			&i.TransactionStatus,
			&i.CurrentUsdtNgnRate,
			&i.CustomerID,
			&i.NgnEquivalent,
			&i.PayoutStatus,
			&i.BankAccName,
			&i.BankAccNumber,
			&i.BitpowrRef,
			&i.TransAddress,
			&i.TransAmount,
			&i.TransChain,
			&i.TransHash,
			&i.BankCode,
			&i.AdminTransAmount,
			&i.AdminTransFee,
			&i.AdminTransRef,
			&i.AdminTransUid,
			&i.TransAmountNgn,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSwapWithBitpowrInfo = `-- name: UpdateSwapWithBitpowrInfo :one
UPDATE coinswap
SET
    bitpowr_ref = COALESCE($1, bitpowr_ref),
    trans_address = COALESCE($2, trans_address),
    trans_amount = COALESCE($3, trans_amount),
    trans_chain = COALESCE($4, trans_chain),
    trans_hash = COALESCE($5, trans_hash),
    transaction_status = COALESCE($6, transaction_status),
    trans_amount_ngn = COALESCE($7, trans_amount_ngn)
WHERE coin_name = $8 AND transaction_status = transaction_status AND coin_address = $9
RETURNING id, coin_name, coin_amount_to_swap, network, phone_number, coin_address, transaction_ref, transaction_status, current_usdt_ngn_rate, customer_id, ngn_equivalent, payout_status, bank_acc_name, bank_acc_number, bitpowr_ref, trans_address, trans_amount, trans_chain, trans_hash, bank_code, admin_trans_amount, admin_trans_fee, admin_trans_ref, admin_trans_uid, trans_amount_ngn, created_at, updated_at
`

type UpdateSwapWithBitpowrInfoParams struct {
	BitpowrRef        sql.NullString `json:"bitpowr_ref"`
	TransAddress      sql.NullString `json:"trans_address"`
	TransAmount       sql.NullString `json:"trans_amount"`
	TransChain        sql.NullString `json:"trans_chain"`
	TransHash         sql.NullString `json:"trans_hash"`
	TransactionStatus string         `json:"transaction_status"`
	TransAmountNgn    sql.NullString `json:"trans_amount_ngn"`
	CoinName          string         `json:"coin_name"`
	CoinAddress       string         `json:"coin_address"`
}

func (q *Queries) UpdateSwapWithBitpowrInfo(ctx context.Context, arg UpdateSwapWithBitpowrInfoParams) (Coinswap, error) {
	row := q.db.QueryRowContext(ctx, updateSwapWithBitpowrInfo,
		arg.BitpowrRef,
		arg.TransAddress,
		arg.TransAmount,
		arg.TransChain,
		arg.TransHash,
		arg.TransactionStatus,
		arg.TransAmountNgn,
		arg.CoinName,
		arg.CoinAddress,
	)
	var i Coinswap
	err := row.Scan(
		&i.ID,
		&i.CoinName,
		&i.CoinAmountToSwap,
		&i.Network,
		&i.PhoneNumber,
		&i.CoinAddress,
		&i.TransactionRef,
		&i.TransactionStatus,
		&i.CurrentUsdtNgnRate,
		&i.CustomerID,
		&i.NgnEquivalent,
		&i.PayoutStatus,
		&i.BankAccName,
		&i.BankAccNumber,
		&i.BitpowrRef,
		&i.TransAddress,
		&i.TransAmount,
		&i.TransChain,
		&i.TransHash,
		&i.BankCode,
		&i.AdminTransAmount,
		&i.AdminTransFee,
		&i.AdminTransRef,
		&i.AdminTransUid,
		&i.TransAmountNgn,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateSwapWithShutterInfo = `-- name: UpdateSwapWithShutterInfo :one
UPDATE  coinswap
SET
   payout_status = COALESCE($1, payout_status),
   admin_trans_amount = COALESCE($2, admin_trans_amount),
   admin_trans_fee = COALESCE($3, admin_trans_fee),
   admin_trans_ref = COALESCE($4, admin_trans_ref)
WHERE transaction_ref = $5
RETURNING id, coin_name, coin_amount_to_swap, network, phone_number, coin_address, transaction_ref, transaction_status, current_usdt_ngn_rate, customer_id, ngn_equivalent, payout_status, bank_acc_name, bank_acc_number, bitpowr_ref, trans_address, trans_amount, trans_chain, trans_hash, bank_code, admin_trans_amount, admin_trans_fee, admin_trans_ref, admin_trans_uid, trans_amount_ngn, created_at, updated_at
`

type UpdateSwapWithShutterInfoParams struct {
	PayoutStatus     sql.NullString `json:"payout_status"`
	AdminTransAmount sql.NullString `json:"admin_trans_amount"`
	AdminTransFee    sql.NullString `json:"admin_trans_fee"`
	AdminTransRef    sql.NullString `json:"admin_trans_ref"`
	TransactionRef   string         `json:"transaction_ref"`
}

func (q *Queries) UpdateSwapWithShutterInfo(ctx context.Context, arg UpdateSwapWithShutterInfoParams) (Coinswap, error) {
	row := q.db.QueryRowContext(ctx, updateSwapWithShutterInfo,
		arg.PayoutStatus,
		arg.AdminTransAmount,
		arg.AdminTransFee,
		arg.AdminTransRef,
		arg.TransactionRef,
	)
	var i Coinswap
	err := row.Scan(
		&i.ID,
		&i.CoinName,
		&i.CoinAmountToSwap,
		&i.Network,
		&i.PhoneNumber,
		&i.CoinAddress,
		&i.TransactionRef,
		&i.TransactionStatus,
		&i.CurrentUsdtNgnRate,
		&i.CustomerID,
		&i.NgnEquivalent,
		&i.PayoutStatus,
		&i.BankAccName,
		&i.BankAccNumber,
		&i.BitpowrRef,
		&i.TransAddress,
		&i.TransAmount,
		&i.TransChain,
		&i.TransHash,
		&i.BankCode,
		&i.AdminTransAmount,
		&i.AdminTransFee,
		&i.AdminTransRef,
		&i.AdminTransUid,
		&i.TransAmountNgn,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
