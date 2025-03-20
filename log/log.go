package log

import (
	"context"

	"github.com/yoshino-s/go-framework/errors"
	"go.uber.org/zap"
)

type loggerKey struct{}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContextOptional(ctx context.Context) (*zap.Logger, bool) {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	return logger, ok
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		zap.L().Warn("logger not found in context")
		return zap.L()
	}
	return logger
}

func FromContextOr(ctx context.Context, fallback *zap.Logger) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		return fallback
	}
	return logger
}

func FromContextX(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		errors.NewMissingComponentError("logger")
	}
	return logger
}
