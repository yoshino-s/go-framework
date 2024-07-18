package mock

import (
	"context"
	"encoding/json"
	"os"

	"github.com/oauth2-proxy/mockoidc"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
	"go.uber.org/zap"
)

type MockOIDCServer struct {
	application.EmptyApplication
	config config
	Server *mockoidc.MockOIDC
}

var _ application.Application = (*MockOIDCServer)(nil)

func New() *MockOIDCServer {
	return &MockOIDCServer{}
}

func (m *MockOIDCServer) Configuration() configuration.Configuration {
	return &m.config
}

func (s *MockOIDCServer) Run(ctx context.Context) {
	s.Logger.Info("Starting mock OIDC server", zap.Any("config", s.config))
	m, err := mockoidc.Run()
	if err != nil {
		panic(err)
	}
	s.Server = m
	if s.config.ClientID != "" {
		m.ClientID = s.config.ClientID
	}
	if s.config.ClientSecret != "" {
		m.ClientSecret = s.config.ClientSecret
	}

	s.Logger.Info("Mock OIDC server started")

	config := m.Config()

	s.Logger.Info("Connect with following config")

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config); err != nil {
		panic(err)
	}

	for range ctx.Done() {
		s.Logger.Info("Mock OIDC server stopped")
	}
}

func (s *MockOIDCServer) Close(context.Context) {
	if err := s.Server.Shutdown(); err != nil {
		panic(err)
	}
}
