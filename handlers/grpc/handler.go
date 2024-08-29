package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type Handler struct {
	*application.EmptyApplication
	*grpc.Server
	config
}

var _ http.Handler = (*Handler)(nil)

func New() *Handler {
	h := &Handler{
		EmptyApplication: application.NewEmptyApplication(),
	}

	return h
}

func (h *Handler) Configuration() configuration.Configuration {
	return &h.config
}

func (h *Handler) Setup(context.Context) {

	unaryServerInterceptor := []grpc.UnaryServerInterceptor{}
	streamServerInterceptor := []grpc.StreamServerInterceptor{}

	unaryServerInterceptor = append(unaryServerInterceptor, logging.UnaryServerInterceptor(interceptorLogger(h.Logger)))
	streamServerInterceptor = append(streamServerInterceptor, logging.StreamServerInterceptor(interceptorLogger(h.Logger)))

	h.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			unaryServerInterceptor...,
		),
		grpc.ChainStreamInterceptor(
			streamServerInterceptor...,
		),
	)
}

func (h *Handler) Run(context.Context) {
	netListener, err := net.Listen("tcp", h.ListenAddr)
	if err != nil {
		panic(err)
	}
	if err := h.Server.Serve(netListener); err != nil {
		panic(err)
	}
}

func (h *Handler) ToClient() (*grpc.ClientConn, error) {
	lis := bufconn.Listen(1024 * 1024)
	go h.Server.Serve(lis)
	conn, err := grpc.NewClient("bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		return nil, err
	}
	return conn, nil
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
