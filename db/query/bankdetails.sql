-- name: GetCustomerBankDetails :one
SELECT * FROM bankdetails WHERE customer_id = $1;

-- name: InsertCustomerBankDetails :many
INSERT INTO bankdetails (
    customer_id,
    bank_name,
    bank_code,
    account_number,
    account_name
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;
