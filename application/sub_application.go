package application

import (
	"context"
	"fmt"
	"time"

	"github.com/sourcegraph/conc/iter"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var _ Application = (*SubApplication)(nil)

type SubApplication struct {
	*EmptyApplication
	sub []Application
}

func NewSubApplication(name string) *SubApplication {
	return &SubApplication{
		EmptyApplication: NewEmptyApplication(name),
		sub:              make([]Application, 0),
	}
}

func (a *SubApplication) Append(sa Application) {
	a.sub = append(a.sub, sa)
}

func (a *SubApplication) Configuration() configuration.Configuration {
	return nil
}

func (a *SubApplication) SetLogger(l *zap.Logger) {
	a.EmptyApplication.SetLogger(l)
	for _, sa := range a.sub {
		sa.SetLogger(l)
	}
}

func (a *SubApplication) Initialize(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageInitialize), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		startTime := time.Now()
		a.Logger.Debug(fmt.Sprintf("initialize sub application %T", *sa), log.Context(ctx))
		(*sa).Initialize(ctx)
		a.Logger.Debug(fmt.Sprintf("sub application %T initialized in %s", *sa, time.Since(startTime)), log.Context(ctx))
	})
}

func (a *SubApplication) Setup(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageSetup), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		startTime := time.Now()
		a.Logger.Debug(fmt.Sprintf("setup sub application %T", *sa), log.Context(ctx))
		(*sa).Setup(ctx)
		a.Logger.Debug(fmt.Sprintf("sub application %T setup in %s", *sa, time.Since(startTime)), log.Context(ctx))
	})
}

func (a *SubApplication) Run(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageRun), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		a.Logger.Debug(fmt.Sprintf("run sub application %T", *sa), log.Context(ctx))
		(*sa).Run(ctx)
	})
}

func (a *SubApplication) Close(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageClose), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		startTime := time.Now()
		a.Logger.Debug(fmt.Sprintf("close sub application %T", *sa), log.Context(ctx))
		(*sa).Close(ctx)
		a.Logger.Debug(fmt.Sprintf("sub application %T closed in %s", *sa, time.Since(startTime)), log.Context(ctx))
	})
}
