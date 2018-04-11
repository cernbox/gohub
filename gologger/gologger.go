package gologger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a configured production *zap.Logger.
func New(level string, outputPaths ...string) *zap.Logger {
	logLevel := zapcore.InfoLevel
	err := logLevel.UnmarshalText([]byte(level))
	if err != nil {
		panic(fmt.Errorf("wrong log level: %s", err.Error()))
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)
	config.OutputPaths = outputPaths

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
