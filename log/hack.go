package log

import (
	"github.com/yoshino-s/go-framework/magic"
	"go.uber.org/zap"
)

func SetLoggerName(logger *zap.Logger, name string) *zap.Logger {
	// change logger.name to name
	if logger.Name() == "" {
		return logger.Named(name)
	}

	clone := *logger

	magic.SetUnexported(logger, "name", name)
	return &clone
}
