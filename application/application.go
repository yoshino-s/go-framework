package application

import (
	"context"

	"github.com/yoshino-s/go-framework/configuration"
	"go.uber.org/zap"
)

type Application interface {
	Configuration() configuration.Configuration
	Setup(context.Context)
	Run(context.Context)
	Close(context.Context)
	SetLogger(*zap.Logger)
}

var _ Application = &MainApplication{}

type EmptyApplication struct {
	Logger *zap.Logger
}

func NewEmptyApplication() *EmptyApplication {
	return &EmptyApplication{
		Logger: zap.NewNop(),
	}
}

func (a *EmptyApplication) Configuration() configuration.Configuration { return nil }
func (a *EmptyApplication) Setup(context.Context)                      {}
func (a *EmptyApplication) Run(context.Context)                        {}
func (a *EmptyApplication) Close(context.Context)                      {}
func (a *EmptyApplication) SetLogger(l *zap.Logger)                    { a.Logger = l }

type FuncApplication func(context.Context)

func (f FuncApplication) Configuration() configuration.Configuration { return nil }
func (f FuncApplication) SetLogger(l *zap.Logger)                    {}
func (f FuncApplication) Setup(ctx context.Context)                  {}
func (f FuncApplication) Run(ctx context.Context)                    { f(ctx) }
func (f FuncApplication) Close(ctx context.Context)                  {}
