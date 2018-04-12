package gologger

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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

// GetLoggedHTTPHandler wraps the provided http.Handlers with Apache logging information.
func GetLoggedHTTPHandler(filename string, h http.Handler) http.Handler {
	var file *os.File
	if filename == "stderr" {
		file = os.Stderr
	} else if filename == "stdout" {
		file = os.Stdout
	} else {

		fd, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("error creating file: %s", err.Error()))
		}
		file = fd
	}

	return handlers.LoggingHandler(file, h)
}
