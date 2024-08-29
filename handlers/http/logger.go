package http

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func toEchoLogger(logger *zap.Logger) echo.Logger {
	return &echoLogger{logger}
}

type echoLogger struct {
	logger *zap.Logger
}

func (l *echoLogger) Output() io.Writer {
	return nil
}

func (l *echoLogger) SetOutput(w io.Writer) {}

func (l *echoLogger) Prefix() string {
	return l.logger.Name()
}

func (l *echoLogger) SetPrefix(p string) {
	l.logger = l.logger.Named(p)
}

func (l *echoLogger) Level() log.Lvl {
	return log.DEBUG
}

func (l *echoLogger) SetLevel(v log.Lvl) {}

func (l *echoLogger) SetHeader(h string) {}

func (l *echoLogger) Print(i ...interface{}) {
	l.logger.Sugar().Info(i...)
}

func (l *echoLogger) Printf(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

func (l *echoLogger) Printj(j log.JSON) {
	l.logger.Sugar().Infow("", j)
}

func (l *echoLogger) Debug(i ...interface{}) {
	l.logger.Sugar().Debug(i...)
}

func (l *echoLogger) Debugf(format string, args ...interface{}) {
	l.logger.Sugar().Debugf(format, args...)
}

func (l *echoLogger) Debugj(j log.JSON) {
	l.logger.Sugar().Debugw("", j)
}

func (l *echoLogger) Info(i ...interface{}) {
	l.logger.Sugar().Info(i...)
}

func (l *echoLogger) Infof(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

func (l *echoLogger) Infoj(j log.JSON) {
	l.logger.Sugar().Infow("", j)
}

func (l *echoLogger) Warn(i ...interface{}) {
	l.logger.Sugar().Warn(i...)
}

func (l *echoLogger) Warnf(format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}

func (l *echoLogger) Warnj(j log.JSON) {
	l.logger.Sugar().Warnw("", j)
}

func (l *echoLogger) Error(i ...interface{}) {
	l.logger.Sugar().Error(i...)
}

func (l *echoLogger) Errorf(format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}

func (l *echoLogger) Errorj(j log.JSON) {
	l.logger.Sugar().Errorw("", j)
}

func (l *echoLogger) Fatal(i ...interface{}) {
	l.logger.Sugar().Fatal(i...)
}

func (l *echoLogger) Fatalj(j log.JSON) {
	l.logger.Sugar().Fatalw("", j)
}

func (l *echoLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}

func (l *echoLogger) Panic(i ...interface{}) {
	l.logger.Sugar().Panic(i...)
}

func (l *echoLogger) Panicj(j log.JSON) {
	l.logger.Sugar().Panicw("", j)
}

func (l *echoLogger) Panicf(format string, args ...interface{}) {
	l.logger.Sugar().Panicf(format, args...)
}
