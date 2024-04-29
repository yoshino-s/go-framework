package application

import (
	"context"
	"fmt"

	"github.com/sourcegraph/conc/iter"
	"github.com/yoshino-s/go-framework/configuration"
	"go.uber.org/zap"
)

var _ Application = (*SubApplication)(nil)

type SubApplication struct {
	*EmptyApplication
	sub []Application
}

func NewSubApplication() *SubApplication {
	return &SubApplication{
		EmptyApplication: NewEmptyApplication(),
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
	a.Logger = l
	for _, sa := range a.sub {
		sa.SetLogger(l)
	}
}

func (a *SubApplication) Setup(ctx context.Context) {
	a.SetLogger(a.Logger)
	iter.ForEach(a.sub, func(sa *Application) {
		a.Logger.Debug("setup sub application", zap.String("application", fmt.Sprintf("%T", *sa)))
		(*sa).Setup(ctx)
	})
}

func (a *SubApplication) Run(ctx context.Context) {
	iter.ForEach(a.sub, func(sa *Application) {
		if *sa != nil {
			a.Logger.Debug("run sub application", zap.String("application", fmt.Sprintf("%T", *sa)))
			(*sa).Run(ctx)
		}
	})
}

func (a *SubApplication) Close(ctx context.Context) {
	iter.ForEach(a.sub, func(sa *Application) {
		a.Logger.Debug("close sub application", zap.String("application", fmt.Sprintf("%T", *sa)))
		(*sa).Close(ctx)
	})
}
