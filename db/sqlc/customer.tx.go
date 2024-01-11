package db

import (
	"context"
)

type CreateCustomerTransactionParams struct {
	CreateCustomerParams
	SendVerificationMail func(customer *Customer) error
}

func (store *Storage) CreateCustomerTransaction(
	ctx context.Context,
	arg CreateCustomerTransactionParams,
) (Customer, error) {
	var customer Customer
	err := store.executeTransaction(ctx, func(q *Queries) error {
		var err error
		customer, err = q.CreateCustomer(ctx, arg.CreateCustomerParams)
		if err != nil {
			return err
		}
		return arg.SendVerificationMail(&customer)
	})
	return customer, err
}
