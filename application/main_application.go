package application

import "context"

type MainApplication struct {
	SubApplication
}

func NewMainApplication() *MainApplication {
	return &MainApplication{
		SubApplication: *NewSubApplication(nil),
	}
}

func (a *MainApplication) Go(ctx context.Context) {
	a.Setup(ctx)
	a.Run(ctx)
	a.Close(ctx)
}
