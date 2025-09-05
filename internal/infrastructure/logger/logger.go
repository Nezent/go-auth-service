package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger provides a production-ready zap.Logger for Uber FX DI.
func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	return cfg.Build()
}

// Helper functions for structured logging with injected logger.
func Info(logger *zap.Logger, msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Error(logger *zap.Logger, msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Debug(logger *zap.Logger, msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Warn(logger *zap.Logger, msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Sync flushes any buffered log entries.
func Sync(logger *zap.Logger) error {
	return logger.Sync()
}
