package logging

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var ZL zerolog.Logger

func New(logLevel string) {
	var lvl zerolog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		lvl = zerolog.DebugLevel
	case "info":
		lvl = zerolog.InfoLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "error":
		lvl = zerolog.ErrorLevel
	case "fatal":
		lvl = zerolog.FatalLevel
	default:
		lvl = zerolog.InfoLevel
	}

	ZL = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		},
	).
		Level(lvl).
		With().Timestamp().Logger()
}

func Debug(msg string) {
	ZL.Debug().Msg(msg)
}

func Debugf(format string, args ...interface{}) {
	ZL.Debug().Msgf(format, args...)
}

func Info(msg string) {
	ZL.Info().Msg(msg)
}

func Infof(format string, args ...interface{}) {
	ZL.Info().Msgf(format, args...)
}

func Warn(msg string) {
	ZL.Warn().Msg(msg)
}

func Warnf(format string, args ...interface{}) {
	ZL.Warn().Msgf(format, args...)
}

func Error(msg string) {
	ZL.Error().Msg(msg)
}

func Errorf(format string, args ...interface{}) {
	ZL.Error().Msgf(format, args...)
}

func ErrorWithErr(err error, msg string) {
	ZL.Error().Err(err).Msg(msg)
}

func ErrorWithErrf(err error, format string, args ...interface{}) {
	ZL.Error().Err(err).Msgf(format, args...)
}

func Fatal(msg string) {
	ZL.Fatal().Msg(msg)
}

func Fatalf(format string, args ...interface{}) {
	ZL.Fatal().Msgf(format, args...)
}
