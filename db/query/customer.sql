-- name: CreateCustomer :one
INSERT INTO customer (
    first_name,
    last_name,
    email,
    password,
    phone
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;
