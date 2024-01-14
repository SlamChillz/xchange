package worker

import (
	"context"
	"encoding/json"
	"database/sql"
	"fmt"

	"github.com/hibiken/asynq"
	db "github.com/slamchillz/xchange/db/sqlc"
	apiTypes "github.com/slamchillz/xchange/types/api"
	"github.com/slamchillz/xchange/utils"
)

const (
	TaskStoreNewCustomer = "task:store_new_customer"
)

type PayloadStoreNewCustomer struct {
	Data []byte `json:"data"`
	TaskOptions []asynq.Option `json:"task_options"`
}

func (TaskDistributor *AsynqTaskDistributor) DistributeTaskStoreNewCustomer(
	ctx context.Context,
	payload PayloadStoreNewCustomer,
) error {
	task := asynq.NewTask(TaskStoreNewCustomer, payload.Data)
	taskInfo, err := TaskDistributor.client.EnqueueContext(ctx, task, payload.TaskOptions...)
	if err != nil {
		return fmt.Errorf("failed to enqueue store-new-customer-in-db task; error: %w", err)
	}
	logger.Info().
		Str("type", taskInfo.Type).
		Bytes("payload", taskInfo.Payload).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msgf("enqueued store-new-customer-in-db task")
	return nil
}

func (processor *AsynqTaskProcessor) ProcessTaskStoreNewCustomer(ctx context.Context, task *asynq.Task) error {
	var payload apiTypes.CreateCustomerRequest
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to json unmarshal store-new-customer-in-db task payload; error: %w", err)
	}
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password in store-new-customer-in-db task; error: %w", err)
	}
	// add logic to send verification email to a new customer
	customer, err := processor.db.CreateCustomer(ctx, db.CreateCustomerParams{
		FirstName: payload.FirstName,
		LastName: payload.LastName,
		Email: payload.Email,
		Password: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		Phone: sql.NullString{
			String: "",
			Valid: false,
		},
		GoogleID: sql.NullString{
			String: "N/A",
			Valid: false,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to store-new-customer-in-db; error: %w", err)
	}
	logger.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("customer-email", customer.Email).
		Msgf("processed store-new-customer-in-db task")
	return nil
}
