-- name: CreateCustomer :one
INSERT INTO customer (
    first_name,
    last_name,
    email,
    password,
    phone,
    photo,
    google_id
) VALUES (
    sqlc.arg('first_name'),
    sqlc.arg('last_name'),
    sqlc.arg('email'),
    sqlc.arg('password'),
    sqlc.arg('phone'),
    sqlc.arg('photo'),
    sqlc.arg('google_id')
) RETURNING *;

-- name: GetCustomerByEmail :one
SELECT * FROM customer WHERE email = $1;

-- name: GetCustomerByPhoneNumber :one
SELECT * FROM customer WHERE phone = $1;

-- name: GetCustomerById :one
SELECT * FROM customer WHERE id = $1;

-- name: GetCustomerByGoogleId :one
SELECT * FROM customer WHERE google_id = $1;

-- name: UpdateCustomerPassword :one
UPDATE customer
SET password = $2
WHERE id = $1
RETURNING *;

-- name: ResetCustomerPassword :one
UPDATE customer
SET password = $2
WHERE email = $1
RETURNING *;
