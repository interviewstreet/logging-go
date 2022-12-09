// Package app Package to initialise and expose application logger
package app

import (
	"context"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/interviewstreet/logging-go/core"
	"github.com/mcuadros/go-defaults"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const TraceIDKey = "trace_id"

var (
	baseLogger    *zap.SugaredLogger
	loggerOnce    sync.Once // Used for ensuring that logger initialisation is only done once.
	contextHeader string
)

func getProcessInfo() map[string]interface{} {
	hostname, _ := os.Hostname()
	return map[string]interface{}{
		"hostname": hostname,
		"pid":      os.Getpid(),
	}
}

func createNewLogger(namespace string, level zapcore.Level, options *core.LoggerOptions) *zap.SugaredLogger {
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
			NameKey:        "logger_name",
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
	return l.Sugar().Named("application")
}

// SetupLogger initialises the logging configuration
func SetupLogger(namespace string, logLevel zapcore.Level, options *core.LoggerOptions) {
	loggerOnce.Do(func() {
		// Populate with defaults
		if options == nil {
			options = &core.LoggerOptions{}
		}
		defaults.SetDefaults(options)
		// Create a new logger
		baseLogger = createNewLogger(namespace, logLevel, options)
		// Set the HTTP header name from where context_id will be extracted
		contextHeader = options.ContextIDHeader
	})
}

// checkInitialisation validates whether the singleton logger variable is properly initialised or not
func checkInitialisation() {
	if baseLogger == nil {
		panic("Logger not initialised, make sure `SetupLogger` is called before")
	}
}

// New returns a logger instance. A new `context_id` is auto-generated during this call.
//
// This function doesn't create a new logger but instead creates a child logger out of already global logger
// initialised during SetupLogger
func New() *zap.SugaredLogger {
	checkInitialisation()
	// Generate a unique context id
	contextID := uuid.NewString()
	return baseLogger.With("context_id", contextID, zap.Namespace("labels"))
}

// NewWithContextID returns a logger instance with `context_id` defined as per the argument.
//
// This function doesn't create a new logger but instead creates a child logger out of already global logger
// initialised during SetupLogger
func NewWithContextID(contextID string) *zap.SugaredLogger {
	checkInitialisation()
	return baseLogger.With("context_id", contextID, zap.Namespace("labels"))
}

// NewWithGinCtx returns a logger instance with gin Context as an argument to extract `context_id` from request
// headers.
// This function doesn't create a new logger but instead creates a child logger out of already global logger
// initialised during SetupLogger
func NewWithGinCtx(ctx *gin.Context) *zap.SugaredLogger {
	checkInitialisation()
	return baseLogger.With("context_id", ctx.GetHeader(contextHeader), zap.Namespace("labels"))
}

// NewWithCtx returns a logger instance with context.Context as an argument to extract `trace_id` from the
// context.
// This function doesn't create a new logger but instead creates a child logger out of already global logger
// initialised during SetupLogger
func NewWithCtx(ctx context.Context) *zap.SugaredLogger {
	checkInitialisation()
	return baseLogger.With("trace_id", ctx.Value(TraceIDKey).(string), zap.Namespace("labels"))
}
