package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/hibiken/asynq"
)

const (
	TaskVerificationEmail = "task:verification_email"
)

type PayloadVerificationEmail struct {
	Email string `json:"email"`
	FirstName string `json:"first_name"`
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
	fmt.Printf("%+v\n", payload)
	msgBuffer := &bytes.Buffer{}
	// add logic to send verification email to a new customer
	mailTemplate, err := template.ParseFiles("templates/verification_email.html")
	if err != nil {
		return fmt.Errorf("failed to parse verification email template: %w", err)
	}
	err = mailTemplate.Execute(msgBuffer, payload)
	if err != nil {
		return fmt.Errorf("failed to execute verification mail template: %w", err)
	}
	fmt.Printf("%s", msgBuffer.Bytes())
	// fmt.Printf("Your OTP is %s\n", payload.Otp)
	logger.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("customer_email", payload.Email).
		Msgf("processed verification email task")
	return nil
}
