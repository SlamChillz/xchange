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
	TaskSendMail = "task:send_mail"
)

type PayloadSendMail struct {
	MailReciever string `json:"mail_reciever"`
	Template string `json:"template"`
	TemplateData interface{} `json:"template_data"`
	MailType string `json:"mail_type"`
	TaskOptions []asynq.Option `json:"task_options"` 
}

func (TaskDistributor *AsynqTaskDistributor) DistributeTaskSendMail(
	ctx context.Context,
	payload PayloadSendMail,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to json marshal send mail task payload; error: %w", err)
	}
	task := asynq.NewTask(TaskSendMail, jsonPayload)
	taskInfo, err := TaskDistributor.client.EnqueueContext(ctx, task, payload.TaskOptions...)
	if err != nil {
		return fmt.Errorf("failed to enqueue send mail task; error: %w", err)
	}
	logger.Info().
		Str("type", taskInfo.Type).
		Bytes("payload", taskInfo.Payload).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Str("mail_type", payload.MailType).
		Str("mail_reciever", payload.MailReciever).
		Msgf("enqueued send mail task")
	return nil
}

func (processor *AsynqTaskProcessor) ProcessTaskSendMail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendMail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to json unmarshal send mail task payload; error: %w", err)
	}
	fmt.Printf("%+v\n", payload)
	msgBuffer := &bytes.Buffer{}
	// add logic to send verification email to a new customer
	mailTemplate, err := template.ParseFiles(payload.Template)
	if err != nil {
		return fmt.Errorf("failed to parse send email template; error: %w", err)
	}
	err = mailTemplate.Execute(msgBuffer, payload.TemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute send mail template; error: %w", err)
	}
	fmt.Printf("%s", msgBuffer.Bytes())
	// fmt.Printf("Your OTP is %s\n", payload.Otp)
	logger.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("mail_type", payload.MailType).
		Str("mail_reciever", payload.MailReciever).
		Msgf("processed send email task")
	return nil
}
