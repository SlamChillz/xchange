-- name: createAdmindata :one
iNSERT INTO admindata (
    bitpowr_account_id,
    btc_address,
    usdt_address,
    usdt_tron_address,
    admin_email
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;
