package log

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		zap.L().Warn("logger not found in context")
		return zap.L()
	}
	return logger
}
