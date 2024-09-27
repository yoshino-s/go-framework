package telemetry

import (
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
)

var _ application.Application = &Telemetry{}

type Telemetry struct {
	*application.EmptyApplication
	config telemetryConfiguration
}

func New() *Telemetry {
	return &Telemetry{
		EmptyApplication: application.NewEmptyApplication(),
	}
}

func (t *Telemetry) Configuration() configuration.Configuration {
	return &t.config
}
