package logger

import (
	// "fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/pkgerrors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/slamchillz/xchange/utils"
)

var (
	once sync.Once
	logger zerolog.Logger
)

// InitLogger initializes the logger
func init() {
	once.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		config, err := utils.LoadConfig("./")
		if err != nil {
			log.Fatal().Err(err).Msg("cannot load configuration file")
		}
		output := os.Stderr
		LOG_LEVEL := os.Getenv("LOG_LEVEL")
		if LOG_LEVEL == "" {
			LOG_LEVEL = "info"
		}
		level, err := zerolog.ParseLevel(config.LOG_LEVEL)
		if err != nil {
			level = zerolog.InfoLevel
		}
		zerolog.SetGlobalLevel(level)
		if config.Env == "dev" {
			logger = zerolog.New(zerolog.ConsoleWriter{Out: output}).With().Caller().Timestamp().Logger()
		} else {
			logger = zerolog.New(output).With().Caller().Timestamp().Logger()
		}
		// output.FormatLevel = func(i interface{}) string {
		// 	return fmt.Sprintf("[%s]", i)
		// }
		// output.FormatMessage = func(i interface{}) string {
		// 	return fmt.Sprintf("| %s |", i)
		// }
		// output.FormatFieldName = func(i interface{}) string {
		// 	return fmt.Sprintf("%s=", i)
		// }
		// output.FormatFieldValue = func(i interface{}) string {
		// 	return fmt.Sprintf("%s", i)
		// }
		// logger = zerolog.New(output).With().Timestamp().Logger()
	})
}

// GetLogger returns the logger
func GetLogger() zerolog.Logger {
	return logger
}
