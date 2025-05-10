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

type ApplicationStage int

const (
	StageBeforeSetup ApplicationStage = iota
	StageSetup
	StageAfterSetup
	StageRun
	StageClose
)

var _ Application = &funcApplication{}

func NewFuncApplication(stage ApplicationStage, f func(context.Context)) *funcApplication {
	return &funcApplication{
		stage: stage,
		f:     f,
	}
}

type funcApplication struct {
	stage ApplicationStage
	f     func(context.Context)
}

func (f *funcApplication) Configuration() configuration.Configuration { return nil }
func (f *funcApplication) SetLogger(l *zap.Logger)                    {}
func (f *funcApplication) BeforeSetup(ctx context.Context) {
	if f.stage == StageBeforeSetup {
		f.f(ctx)
	}
}
func (f *funcApplication) Setup(ctx context.Context) {
	if f.stage == StageSetup {
		f.f(ctx)
	}
}
func (f *funcApplication) AfterSetup(ctx context.Context) {
	if f.stage == StageAfterSetup {
		f.f(ctx)
	}
}
func (f *funcApplication) Run(ctx context.Context) {
	if f.stage == StageRun {
		f.f(ctx)
	}
}
func (f *funcApplication) Close(ctx context.Context) {
	if f.stage == StageClose {
		f.f(ctx)
	}
}
