// Package to initialise and expose application logger
package app

import (
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/interviewstreet/logging-go/core"
	"github.com/mcuadros/go-defaults"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	baseLogger *zap.SugaredLogger
	loggerOnce sync.Once // Used for ensuring that logger initialisation is only done once.
)

func getProcessInfo() map[string]interface{} {
	hostname, _ := os.Hostname()
	return map[string]interface{}{
		"hostname": hostname,
		"pid":      os.Getpid(),
	}
}

func createNewLogger(namespace string, level zapcore.Level, options *core.LoggerOptions) *zap.SugaredLogger {
	if options == nil {
		options = &core.LoggerOptions{}
	}
	defaults.SetDefaults(options)

	atom := zap.NewAtomicLevel()
	atom.SetLevel(level)

	config := zap.Config{
		Level:            atom,
		Encoding:         "json",
		OutputPaths:      []string{options.OutputPath},
		ErrorOutputPaths: []string{options.ErrorOutputPath},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "severity",
			NameKey:        zapcore.OmitKey,
			FunctionKey:    "source_function",
			CallerKey:      "source_caller",
			MessageKey:     "text_payload",
			StacktraceKey:  "error_stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.EpochMillisTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		InitialFields: map[string]interface{}{
			"namespace":       namespace,
			"environment":     options.Env,
			"resource_labels": getProcessInfo(),
		},
	}
	l, _ := config.Build(zap.AddStacktrace(zapcore.ErrorLevel), zap.Fields())
	return l.Sugar()
}

// SetupLogger initialises the logging configuration
func SetupLogger(namespace string, logLevel zapcore.Level, options *core.LoggerOptions) {
	loggerOnce.Do(func() {
		baseLogger = createNewLogger(namespace, logLevel, options)
	})
}

// NewWithCtx returns a logger instance with `context_id` defined as per the argument.
//
// This function doesn't create a new logger but instead creates a child logger out of already global logger
// initialised during SetupLogger
func NewWithCtx(contextID string) *zap.SugaredLogger {
	return baseLogger.With("context_id", contextID, zap.Namespace("labels"))
}

// New returns a logger instance. A new `context_id` is auto-generated during this call.
//
// This function doesn't create a new logger but instead creates a child logger out of already global logger
// initialised during SetupLogger
func New() *zap.SugaredLogger {
	// Generate a unique context id
	contextID := uuid.NewString()
	return baseLogger.With("context_id", contextID, zap.Namespace("labels"))
}
