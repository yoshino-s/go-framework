package grpc_gateway

import (
	"net/http"

	grpcRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type GHandlerOptions struct {
}

var HandlerOptions = GHandlerOptions{}

type Handler struct {
	*grpcRuntime.ServeMux
}

var _ http.Handler = (*Handler)(nil)

func New() (*Handler, error) {
	h := &Handler{}

	grpcGatewayMux := grpcRuntime.NewServeMux()
	h.ServeMux = grpcGatewayMux

	return h, nil
}
