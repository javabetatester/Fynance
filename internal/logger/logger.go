package logger

import (
	"Fynance/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var globalLogger zerolog.Logger

func Init(cfg *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = false

	level := parseLogLevel(cfg.App.LogLevel)

	zerolog.SetGlobalLevel(level)

	var output zerolog.ConsoleWriter
	if cfg.App.Environment == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
	} else {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		}
	}

	globalLogger = zerolog.New(output).
		With().
		Timestamp().
		Str("service", "fynance").
		Str("environment", cfg.App.Environment).
		Logger()

	log.Logger = globalLogger
}

func GetLogger() zerolog.Logger {
	return globalLogger
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func Debug() *zerolog.Event {
	return globalLogger.Debug()
}

func Info() *zerolog.Event {
	return globalLogger.Info()
}

func Warn() *zerolog.Event {
	return globalLogger.Warn()
}

func Error() *zerolog.Event {
	return globalLogger.Error()
}

func Fatal() *zerolog.Event {
	return globalLogger.Fatal()
}

func Panic() *zerolog.Event {
	return globalLogger.Panic()
}

func With() zerolog.Context {
	return globalLogger.With()
}
