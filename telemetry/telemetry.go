package telemetry

import (
	"context"

	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
)

var _ application.Application = &Telemetry{}

type Telemetry struct {
	*application.EmptyApplication
	config config
}

func New() *Telemetry {
	return &Telemetry{
		EmptyApplication: application.NewEmptyApplication(),
		config:           config{},
	}
}

func (t *Telemetry) Configuration() configuration.Configuration {
	return &telemetryConfiguration{
		config: &t.config,
	}
}

func (t *Telemetry) Setup(context.Context) {
	if t.config.SentryDSN != "" {
		initSentry(t.config.SentryDSN)
	}
}
