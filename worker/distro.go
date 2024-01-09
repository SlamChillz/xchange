package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	// DistributeTask(ctx context.Context, task Task) error
	DistributeTaskVerificationEmail(
		ctx context.Context,
		payload PayloadVerificationEmail,
		opts ...asynq.Option,
	) error
}

type AsynqTaskDistributor struct {
	client *asynq.Client
}

func NewAsynqTaskDistributor(asynqClientOpt asynq.RedisClientOpt) TaskDistributor {
	asynqClient := asynq.NewClient(asynqClientOpt)
	return &AsynqTaskDistributor{client: asynqClient}
}
