package application

import (
	"context"

	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/log"
	"go.uber.org/zap"
)

type Application interface {
	Configuration() configuration.Configuration
	BeforeSetup(context.Context)
	Setup(context.Context)
	AfterSetup(context.Context)
	Run(context.Context)
	Close(context.Context)
	SetLogger(*zap.Logger)
}

var _ Application = &EmptyApplication{}

type EmptyApplication struct {
	Logger *zap.Logger
	Name   string
}

func NewEmptyApplication(name string) *EmptyApplication {
	return &EmptyApplication{
		Logger: zap.NewNop(),
		Name:   name,
	}
}

func (a *EmptyApplication) Configuration() configuration.Configuration { return nil }
func (a *EmptyApplication) BeforeSetup(context.Context)                {}
func (a *EmptyApplication) Setup(context.Context)                      {}
func (a *EmptyApplication) AfterSetup(context.Context)                 {}
func (a *EmptyApplication) Run(context.Context)                        {}
func (a *EmptyApplication) Close(context.Context)                      {}
func (a *EmptyApplication) SetLogger(l *zap.Logger) {
	a.Logger = log.SetLoggerName(l, a.Name)
}

var _ Application = FuncApplication(nil)

type FuncApplication func(context.Context)

func (f FuncApplication) Configuration() configuration.Configuration { return nil }
func (f FuncApplication) SetLogger(l *zap.Logger)                    {}
func (f FuncApplication) BeforeSetup(ctx context.Context)            {}
func (f FuncApplication) Setup(ctx context.Context)                  {}
func (f FuncApplication) AfterSetup(ctx context.Context)             {}
func (f FuncApplication) Run(ctx context.Context)                    { f(ctx) }
func (f FuncApplication) Close(ctx context.Context)                  {}
