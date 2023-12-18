// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: customer.sql

package db

import (
	"context"
)

const createCustomer = `-- name: CreateCustomer :one
INSERT INTO customer (
    first_name,
    last_name,
    email,
    password,
    phone
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, last_login, photo, first_name, last_name, email, password, phone, is_active, is_staff, is_supercustomer, created_at, updated_at
`

type CreateCustomerParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

func (q *Queries) CreateCustomer(ctx context.Context, arg CreateCustomerParams) (Customer, error) {
	row := q.db.QueryRowContext(ctx, createCustomer,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
		arg.Phone,
	)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.LastLogin,
		&i.Photo,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Phone,
		&i.IsActive,
		&i.IsStaff,
		&i.IsSupercustomer,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCustomerByEmail = `-- name: GetCustomerByEmail :one
SELECT id, last_login, photo, first_name, last_name, email, password, phone, is_active, is_staff, is_supercustomer, created_at, updated_at FROM customer WHERE email = $1
`

func (q *Queries) GetCustomerByEmail(ctx context.Context, email string) (Customer, error) {
	row := q.db.QueryRowContext(ctx, getCustomerByEmail, email)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.LastLogin,
		&i.Photo,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Phone,
		&i.IsActive,
		&i.IsStaff,
		&i.IsSupercustomer,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCustomerByPhoneNumber = `-- name: GetCustomerByPhoneNumber :one
SELECT id, last_login, photo, first_name, last_name, email, password, phone, is_active, is_staff, is_supercustomer, created_at, updated_at FROM customer WHERE phone = $1
`

func (q *Queries) GetCustomerByPhoneNumber(ctx context.Context, phone string) (Customer, error) {
	row := q.db.QueryRowContext(ctx, getCustomerByPhoneNumber, phone)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.LastLogin,
		&i.Photo,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Phone,
		&i.IsActive,
		&i.IsStaff,
		&i.IsSupercustomer,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
