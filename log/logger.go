package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
)

func initialize(instance *zap.Logger) {
	if defaultLogger != nil {
		defaultLogger.Debug("Replacing logger factory", zap.String("old", defaultLogger.Name()), zap.String("new", instance.Name()))
	} else {
		instance.Debug("Initializing logger factory", zap.String("factory", instance.Name()))
	}
	defaultLogger = instance
}

func init() {
	initialize(zap.L())
}


// Debug logs an debug msg with fields
func Debug(msg string, fields ...zapcore.Field) {
	defaultLogger.Debug(msg, fields...)
}

// Info logs an info msg with fields
func Info(msg string, fields ...zapcore.Field) {
	defaultLogger.Info(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func Fatal(msg string, fields ...zapcore.Field) {
	defaultLogger.Fatal(msg, fields...)
}

