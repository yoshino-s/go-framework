package application

import "go.uber.org/zap/zapcore"

type LoggerCoreContrib interface {
	LoggerCores() []zapcore.Core
}

func getLoggerCores(app Application) []zapcore.Core {
	if app, ok := app.(LoggerCoreContrib); ok {
		return app.LoggerCores()
	}
	return nil
}
