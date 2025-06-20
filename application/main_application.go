package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sourcegraph/conc/iter"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/log"
	"go.uber.org/zap"
)

var _ Application = (*MainApplication)(nil)

type MainApplication struct {
	*SubApplication
	*Container
	signalChannel chan os.Signal
}

func NewMainApplication() *MainApplication {
	return &MainApplication{
		SubApplication: NewSubApplication("MainApplication"),
		Container:      &Container{},
	}
}

func (a *MainApplication) Configuration() configuration.Configuration {
	return &configuration.CombinationConfiguration{
		&log.LogConfiguration{
			Logger: &a.Logger,
		},
	}
}

func (a *MainApplication) Setup(ctx context.Context) {
	a.signalChannel = make(chan os.Signal)
	signal.Notify(a.signalChannel, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)

	a.SubApplication.Setup(ctx)
}

func (a *MainApplication) Go(ctx context.Context) {
	a.SetLogger(a.Logger)

	a.Initialize(ctx)

	iter.ForEach(a.sub, func(sa *Application) {
		if err := a.doSet(*sa); err != nil {
			a.Logger.Fatal("Failed to set application", zap.Error(err))
		}
	})

	a.Setup(ctx)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for v := range a.signalChannel {
			switch v {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT:
				a.Logger.Debug("Receive signal", zap.Any("signal", v))
				cancel()
			default:
				a.Logger.Debug("Receive unknown signal", zap.Any("signal", v))
			}
		}
	}()
	go func() {
		a.Run(ctx)
		cancel()
	}()

	<-ctx.Done()

	a.Logger.Debug("Close MainApplication")
	a.Close(ctx)
	a.Logger.Debug("Bye!")
}

func (a *MainApplication) Append(sa Application) {
	a.SubApplication.Append(sa)
	a.Register(sa)
}
