-- name: GetBtcAddress :one
SELECT btc_address FROM customerasset WHERE customer_id = $1;

-- name: GetUsdtAddress :one
SELECT usdt_address FROM customerasset WHERE customer_id = $1;

-- name: GetUsdtTronAddress :one
SELECT usdt_tron_address FROM customerasset WHERE customer_id = $1;

-- name: GetUsdtBscAddress :one
SELECT usdt_bsc_address FROM customerasset WHERE customer_id = $1;

-- name: InsertNewBtcAddress :one
INSERT INTO customerasset (btc_network, btc_address, btc_address_uid, customer_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: InsertNewUsdtAddress :one
INSERT INTO customerasset (usdt_network, usdt_address, usdt_address_uid, customer_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: InsertNewUsdtTronAddress :one
INSERT INTO customerasset (usdt_tron_network, usdt_tron_address, usdt_tron_address_uid, customer_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: InsertNewUsdtBscAddress :one
INSERT INTO customerasset (usdt_bsc_network, usdt_bsc_address, usdt_bsc_address_uid, customer_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateBtcAddress :one
UPDATE customerasset SET btc_network = $1, btc_address = $2, btc_address_uid = $3 WHERE customer_id = $4 RETURNING *;

-- name: UpdateUsdtAddress :one
UPDATE customerasset SET usdt_network = $1, usdt_address = $2, usdt_address_uid = $3 WHERE customer_id = $4 RETURNING *;

-- name: UpdateUsdtTronAddress :one
UPDATE customerasset SET usdt_tron_network = $1, usdt_tron_address = $2, usdt_tron_address_uid = $3 WHERE customer_id = $4 RETURNING *;

-- name: UpdateUsdtBscAddress :one
UPDATE customerasset SET usdt_bsc_network = $1, usdt_bsc_address = $2, usdt_bsc_address_uid = $3 WHERE customer_id = $4 RETURNING *;
