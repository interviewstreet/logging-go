// Package to initialise and expose middleware/s which support structured request logging
package request

import (
	"fmt"
	"os"
	"sync"

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

func createNewLogger(namespace string, options *core.LoggerOptions) *zap.SugaredLogger {
	if options == nil {
		options = &core.LoggerOptions{}
	}
	defaults.SetDefaults(options)

	atom := zap.NewAtomicLevel()
	atom.SetLevel(zap.InfoLevel)

	config := zap.Config{
		Level:       atom,
		Encoding:    "json",
		OutputPaths: []string{options.OutputPath},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:    "timestamp",
			LevelKey:   zapcore.OmitKey,
			MessageKey: zapcore.OmitKey,
			NameKey:    "request_logger",
			LineEnding: zapcore.DefaultLineEnding,
			EncodeTime: zapcore.EpochMillisTimeEncoder,
		},
		InitialFields: map[string]interface{}{
			"namespace":       namespace,
			"environment":     options.Env,
			"resource_labels": getProcessInfo(),
		},
	}
	l, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialise request logger, error: %s", err))
	}
	return l.Sugar()
}
