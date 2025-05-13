package application

import (
	"context"
	"fmt"

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

func (a *SubApplication) BeforeSetup(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageBeforeSetup), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		a.Logger.Debug(fmt.Sprintf("before setup sub application %T", *sa), log.Context(ctx))
		(*sa).BeforeSetup(ctx)
	})
}

func (a *SubApplication) Setup(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageSetup), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		a.Logger.Debug(fmt.Sprintf("setup sub application %T", *sa), log.Context(ctx))
		(*sa).Setup(ctx)
	})
}

func (a *SubApplication) AfterSetup(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageAfterSetup), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		a.Logger.Debug(fmt.Sprintf("after setup sub application %T", *sa), log.Context(ctx))
		(*sa).AfterSetup(ctx)
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
		a.Logger.Debug(fmt.Sprintf("close sub application %T", *sa), log.Context(ctx))
		(*sa).Close(ctx)
	})
}
