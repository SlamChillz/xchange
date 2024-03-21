package worker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	log "github.com/slamchillz/xchange/logger"
)

var (
	logger = log.GetLogger()
)

type TaskLogger struct {}

func NewTaskLogger() *TaskLogger {
	return &TaskLogger{}
}

func (tl *TaskLogger) Print(level zerolog.Level, args ...interface{}) {
	logger.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (tl *TaskLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	logger.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (tl *TaskLogger) Debug(args ...interface{}) {
	tl.Print(zerolog.DebugLevel, args...)
}

func (tl *TaskLogger) Info(args ...interface{}) {
	tl.Print(zerolog.InfoLevel, args...)
}

func (tl *TaskLogger) Warn(args ...interface{}) {
	tl.Print(zerolog.WarnLevel, args...)
}

func (tl *TaskLogger) Error(args ...interface{}) {
	tl.Print(zerolog.ErrorLevel, args...)
}

func (tl *TaskLogger) Fatal(args ...interface{}) {
	tl.Print(zerolog.FatalLevel, args...)
}
