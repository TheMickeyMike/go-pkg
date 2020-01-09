package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
)

func initialize(instance *zap.Logger) {
	defaultLogger = instance
}

func init() {
	initialize(zap.L())
}

// Info logs an info msg with fields
func Info(msg string, fields ...zapcore.Field) {
	defaultLogger.Info(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func Fatal(msg string, fields ...zapcore.Field) {
	defaultLogger.Fatal(msg, fields...)
}
