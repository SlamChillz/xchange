// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: admindata.sql

package db

import (
	"context"
	"database/sql"
)

const createAdmindata = `-- name: createAdmindata :one
iNSERT INTO admindata (
    bitpowr_account_id,
    btc_address,
    usdt_address,
    usdt_tron_address,
    admin_email
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, bitpowr_account_id, btc_address, usdt_address, usdt_tron_address, admin_email, created_at, updated_at
`

type createAdmindataParams struct {
	BitpowrAccountID string         `json:"bitpowr_account_id"`
	BtcAddress       sql.NullString `json:"btc_address"`
	UsdtAddress      sql.NullString `json:"usdt_address"`
	UsdtTronAddress  sql.NullString `json:"usdt_tron_address"`
	AdminEmail       string         `json:"admin_email"`
}

func (q *Queries) createAdmindata(ctx context.Context, arg createAdmindataParams) (Admindatum, error) {
	row := q.db.QueryRowContext(ctx, createAdmindata,
		arg.BitpowrAccountID,
		arg.BtcAddress,
		arg.UsdtAddress,
		arg.UsdtTronAddress,
		arg.AdminEmail,
	)
	var i Admindatum
	err := row.Scan(
		&i.ID,
		&i.BitpowrAccountID,
		&i.BtcAddress,
		&i.UsdtAddress,
		&i.UsdtTronAddress,
		&i.AdminEmail,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
