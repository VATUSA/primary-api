package logger

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int8

const (
	Debug LogLevel = iota - 1
	Info
	Warn
)

var Logger *zap.Logger

func Setup() {
	level := config.Cfg.API.LoggerLevel
	logLevel := Debug
	switch level {
	case "debug":
		logLevel = Debug
	case "info":
		logLevel = Info
	case "warn":
		logLevel = Warn
	}

	var config zap.Config
	if logLevel == Debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	switch logLevel {
	case Debug:
		config.Level.SetLevel(zap.DebugLevel)
	case Info:
		config.Level.SetLevel(zap.InfoLevel)
	case Warn:
		config.Level.SetLevel(zap.WarnLevel)
	}

	config.EncoderConfig.CallerKey = ""
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
}
