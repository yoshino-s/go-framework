package application

import (
	"context"
	"fmt"

	"github.com/sourcegraph/conc/iter"
	"github.com/yoshino-s/go-framework/configuration"
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

	a.Logger.Debug("before setup sub application", zap.String("application", fmt.Sprintf("%T", *a)))
	iter.ForEach(a.sub, func(sa *Application) {
		if *sa != nil {
			a.Logger.Debug("before setup sub application", zap.String("application", fmt.Sprintf("%T", *sa)))
			(*sa).BeforeSetup(ctx)
		}
	})
}

func (a *SubApplication) Setup(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageSetup), trace.WithNewRoot())
	defer span.End()

	a.Logger.Debug("setup sub application", zap.String("application", fmt.Sprintf("%T", *a)))
	iter.ForEach(a.sub, func(sa *Application) {
		a.Logger.Debug("setup sub application", zap.String("application", fmt.Sprintf("%T", *sa)))
		(*sa).Setup(ctx)
	})
}

func (a *SubApplication) AfterSetup(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageAfterSetup), trace.WithNewRoot())
	defer span.End()

	a.Logger.Debug("after setup sub application", zap.String("application", fmt.Sprintf("%T", *a)))
	iter.ForEach(a.sub, func(sa *Application) {
		if *sa != nil {
			a.Logger.Debug("after setup sub application", zap.String("application", fmt.Sprintf("%T", *sa)))
			(*sa).AfterSetup(ctx)
		}
	})
}

func (a *SubApplication) Run(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageRun), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		if *sa != nil {
			a.Logger.Debug("run sub application", zap.String("application", fmt.Sprintf("%T", *sa)))
			(*sa).Run(ctx)
		}
	})
}

func (a *SubApplication) Close(ctx context.Context) {
	ctx, span := otel.Tracer(ScopeName).Start(ctx, a.spanName(StageClose), trace.WithNewRoot())
	defer span.End()

	iter.ForEach(a.sub, func(sa *Application) {
		a.Logger.Debug("close sub application", zap.String("application", fmt.Sprintf("%T", *sa)))
		(*sa).Close(ctx)
	})
}
