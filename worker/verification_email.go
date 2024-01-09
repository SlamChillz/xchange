package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	db "github.com/slamchillz/xchange/db/sqlc"
)

const (
	TaskVerificationEmail = "task:verification_email"
)

type PayloadVerificationEmail struct {
	Customer *db.Customer `json:"customer"`
	Otp string `json:"otp"`
}

func (TaskDistributor *AsynqTaskDistributor) DistributeTaskVerificationEmail(
	ctx context.Context,
	payload PayloadVerificationEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to json marshal verification email task payload: %w", err)
	}
	task := asynq.NewTask(TaskVerificationEmail, jsonPayload)
	taskInfo, err := TaskDistributor.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue verification email task: %w", err)
	}
	logger.Info().
		Str("type", taskInfo.Type).
		Bytes("payload", taskInfo.Payload).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msgf("enqueued verification email task")
	return nil
}

func (processor *AsynqTaskProcessor) ProcessTaskVerificationEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadVerificationEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to json unmarshal verification email task payload: %w", err)
	}
	// add logic to send verification email to a new customer
	fmt.Printf("Your OTP is %s\n", payload.Otp)
	logger.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("customer_email", payload.Customer.Email).
		Msgf("processed verification email task")
	return nil
}
