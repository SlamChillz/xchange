package db

import (
	"context"
	"testing"

	"github.com/slamchillz/xchange/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCustomer(t *testing.T) Customer {
	arg := CreateCustomerParams{
		FirstName: utils.RandomName(),
		LastName:  utils.RandomName(),
		Email:     utils.RandomEmail(),
		Password:  utils.RandomString(8),
		Phone:     utils.RandomPhoneNumber(),
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
