package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options declares logger options for builder
type Options struct {
	Debug    bool
	LogLevel string
	AppName  string
	AppID    string
	Version  string
	Revision string
}

// DefaultOptions defines default logger options
var DefaultOptions = &Options{
	Debug:    false,
	LogLevel: "info",
	AppName:  "changeme",
	AppID:    "changeme",
	Version:  "0.0.1",
	Revision: "123456789",
}

func Setup(ctx context.Context, opts Options) {

	// Initialize logs
	var config zap.Config

	if opts.Debug {
		config = zap.NewDevelopmentConfig()
		config.DisableCaller = true
		config.DisableStacktrace = true
	} else {
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
		config.EncoderConfig.MessageKey = "@message"
		config.EncoderConfig.TimeKey = "@timestamp"
		config.EncoderConfig.CallerKey = "@caller"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Parse log level
	errLogLevel := config.Level.UnmarshalText([]byte(opts.LogLevel))
	if errLogLevel != nil {
		panic(errLogLevel)
	}

	// Build real logger
	logger, err := config.Build(
		zap.AddCallerSkip(2),
	)
	if err != nil {
		panic(err)
	}

	// Add prefix to logger
	logger = logger.With(
		zap.String("@appName", opts.AppName),
		zap.String("@version", opts.Version),
		zap.String("@revision", opts.Revision),
		zap.String("@appID", opts.AppID),
		zap.Namespace("@fields"),
	)

	// Override the global factory
	initialize(logger)

	// Override zap default logger
	zap.ReplaceGlobals(logger)
}
