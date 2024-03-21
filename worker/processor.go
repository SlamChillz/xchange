package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/redisdb"
)

var (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	// DistributeTask(ctx context.Context, task Task) error
	Start() error
	ProcessTaskSendMail(
		ctx context.Context,
		task *asynq.Task,
	) error
	ProcessTaskStoreNewCustomer(
		ctx context.Context,
		task *asynq.Task,
	) error
}

type AsynqTaskProcessor struct {
	server *asynq.Server
	db db.Store
	redisdb redisdb.RedisClient
}

func NewAsynqTaskProcessor(asynqServerOpt asynq.RedisClientOpt, db db.Store, redisdb redisdb.RedisClient) TaskProcessor {
	taskLogger := NewTaskLogger()
	redis.SetLogger(taskLogger)

	asynqServer := asynq.NewServer(
		asynqServerOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				// log error messages to log file
				logger.Error().Err(err).
					Str("task", task.Type()).
					Bytes("payload", task.Payload()).
					Msgf("error processing %s task\n", task.Type())
			}),
			Logger: taskLogger,
		},
	)
	return &AsynqTaskProcessor{server: asynqServer, db: db, redisdb: redisdb}
}

func (atp *AsynqTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendMail, atp.ProcessTaskSendMail)
	mux.HandleFunc(TaskStoreNewCustomer, atp.ProcessTaskStoreNewCustomer)
	return atp.server.Start(mux)
}
