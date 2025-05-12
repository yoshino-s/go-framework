package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Context(ctx context.Context) zap.Field {
	return zap.Field{
		Type:      zapcore.SkipType,
		Interface: ctx,
	}
}
