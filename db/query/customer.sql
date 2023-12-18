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

-- name: GetCustomerByEmail :one
SELECT * FROM customer WHERE email = $1;

-- name: GetCustomerByPhoneNumber :one
SELECT * FROM customer WHERE phone = $1;
