package app

import (
	"context"

	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/configuration"
)

var _ application.Application = (*DemoApp)(nil)

type DemoApp struct {
	*application.EmptyApplication
}

func New() *DemoApp {
	return &DemoApp{
		application.NewEmptyApplication(),
	}
}

func (a *DemoApp) Configuration() configuration.Configuration {
	return nil
}

func (a *DemoApp) Setup(context.Context) {
	a.Logger.Info("setup")
}
func (a *DemoApp) Run(context.Context) {
	a.Logger.Info("run")
}
func (a *DemoApp) Close(context.Context) {
	a.Logger.Info("close")
}
