package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yoshino-s/go-framework/configuration"
	"go.uber.org/zap"
)

var _ Application = (*MainApplication)(nil)

type MainApplication struct {
	*SubApplication
	signalChannel chan os.Signal
}

func NewMainApplication() *MainApplication {
	return &MainApplication{
		SubApplication: NewSubApplication(),
	}
}

func (a *MainApplication) Configuration() configuration.Configuration {
	return &configuration.CombinationConfiguration{
		&logConfiguration{
			logger: &a.Logger,
		},
	}
}

func (a *MainApplication) Setup(ctx context.Context) {
	a.signalChannel = make(chan os.Signal)
	signal.Notify(a.signalChannel, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)

	a.SubApplication.Setup(ctx)
}

func (a *MainApplication) Go(ctx context.Context) {
	a.Logger.Debug("Setup MainApplication")
	a.Setup(ctx)
	a.Logger.Debug("Run MainApplication")

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
