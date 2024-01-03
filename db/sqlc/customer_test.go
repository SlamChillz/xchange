package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCustomer(t *testing.T) Customer {
	// hashedPassword := utils.RandomString(8)
	arg := CreateCustomerParams{
		FirstName: utils.RandomName(),
		LastName:  utils.RandomName(),
		Email:     utils.RandomEmail(),
		Password:  sql.NullString{
			String: utils.RandomString(8),
			Valid:  true,
		},
		Phone: sql.NullString{
			String: utils.RandomPhoneNumber(),
			Valid:  true,
		},
	}
	customer, err := testQueries.CreateCustomer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, customer)
	require.Equal(t, arg.FirstName, customer.FirstName)
	require.Equal(t, arg.LastName, customer.LastName)
	require.Equal(t, arg.Email, customer.Email)
	require.Equal(t, arg.Password, customer.Password)
	require.Equal(t, arg.Phone, customer.Phone)
	require.NotZero(t, customer.ID)
	require.NotZero(t, customer.CreatedAt)
	require.NotZero(t, customer.UpdatedAt)
	return customer
}

func TestCreateCustomer(t *testing.T) {
	createRandomCustomer(t)
}
