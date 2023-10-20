-- name: RecentUsdtNgnRate :one
INSERT INTO usdtngnrate (
    usdt_ngn_rate
) VALUES (
    $1
) RETURNING *;
