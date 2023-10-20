-- name: CreateSwap :one
INSERT INTO coinswap (
    coin_name, coin_amount_to_swap, network, phone_number, coin_address, transaction_ref,
    transaction_status, current_usdt_ngn_rate, customer_id, ngn_equivalent,
    payout_status, bank_acc_name, bank_acc_number, bank_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $13, $14
) RETURNING *;


-- name: UpdateSwap :one
UPDATE coinswap 
SET bitpowr_ref = $1,
    trans_address = $2,
    trans_amount = $3,
    trans_chain = $4,
    trans_hash = $5,
    admin_trans_amount = $6,
    admin_trans_fee = $7,
    admin_trans_ref = $8,
    admin_trans_uid = $9,
    trans_amount_ngn = $10
WHERE id = $11
RETURNING *;
