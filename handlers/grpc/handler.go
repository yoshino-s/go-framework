package grpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Handler struct {
	*grpc.Server
}

var _ http.Handler = (*Handler)(nil)

func New(config Config) (*Handler, error) {
	h := &Handler{}

	unaryServerInterceptor := []grpc.UnaryServerInterceptor{}
	streamServerInterceptor := []grpc.StreamServerInterceptor{}

	if config.Log {
		logging.UnaryServerInterceptor(interceptorLogger(zap.L()))
	}

	h.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			unaryServerInterceptor...,
		),
		grpc.ChainStreamInterceptor(
			streamServerInterceptor...,
		),
	)

	return h, nil
}

func interceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		zapFields := make([]zap.Field, 0, len(fields)/2)
		for i := 0; i < len(fields); i += 2 {
			zapFields = append(zapFields, zap.Any(fields[i].(string), fields[i+1]))
		}
		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, zapFields...)
		case logging.LevelInfo:
			l.Info(msg, zapFields...)
		case logging.LevelWarn:
			l.Warn(msg, zapFields...)
		case logging.LevelError:
			l.Error(msg, zapFields...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
