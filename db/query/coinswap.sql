-- name: CreateSwap :one
INSERT INTO coinswap (
    coin_name, coin_amount_to_swap, network, phone_number,
    coin_address, transaction_ref, transaction_status, current_usdt_ngn_rate,
    customer_id, bank_acc_name, bank_acc_number, bank_code
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: UpdateSwapWithBitpowrInfo :one
UPDATE coinswap
SET
    bitpowr_ref = COALESCE(sqlc.narg('bitpowr_ref'), bitpowr_ref),
    trans_address = COALESCE(sqlc.narg('trans_address'), trans_address),
    trans_amount = COALESCE(sqlc.narg('trans_amount'), trans_amount),
    trans_chain = COALESCE(sqlc.narg('trans_chain'), trans_chain),
    trans_hash = COALESCE(sqlc.narg('trans_hash'), trans_hash),
    transaction_status = COALESCE(sqlc.arg('transaction_status'), transaction_status),
    trans_amount_ngn = COALESCE(sqlc.narg('trans_amount_ngn'), trans_amount_ngn)
WHERE coin_name = sqlc.arg('coin_name') AND transaction_status = transaction_status AND coin_address = sqlc.arg('coin_address')
RETURNING *;

-- name: UpdateSwapWithShutterInfo :one
UPDATE  coinswap
SET
   payout_status = COALESCE(sqlc.narg('payout_status'), payout_status),
   admin_trans_amount = COALESCE(sqlc.narg('admin_trans_amount'), admin_trans_amount),
   admin_trans_fee = COALESCE(sqlc.narg('admin_trans_fee'), admin_trans_fee),
   admin_trans_ref = COALESCE(sqlc.narg('admin_trans_ref'), admin_trans_ref)
WHERE transaction_ref = sqlc.arg('transaction_ref')
RETURNING *;

-- name: GetPendingNetworkTransaction :one
SELECT COUNT(*) FROM coinswap WHERE customer_id = $1 AND network = $2 AND transaction_status = $3;
