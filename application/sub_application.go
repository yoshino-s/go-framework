package application

import (
	"context"

	"github.com/sourcegraph/conc"
)

type SubApplication struct {
	application     Application
	subApplications []Application
}

func NewSubApplication(a Application) *SubApplication {
	return &SubApplication{
		application:     a,
		subApplications: []Application{},
	}
}

func (a *SubApplication) Add(sa Application) {
	a.subApplications = append(a.subApplications, sa)
}

func (a *SubApplication) Setup(ctx context.Context) { // sub first, then main
	wg := conc.NewWaitGroup()
	for _, sa := range a.subApplications {
		func(sa Application) {
			wg.Go(func() {
				sa.Setup(ctx)
			})
		}(sa)
	}
	wg.Wait()

	if a.application != nil {
		a.application.Setup(ctx)
	}
}

func (a *SubApplication) Run(ctx context.Context) { // main and sub concurrently
	wg := conc.NewWaitGroup()
	for _, sa := range a.subApplications {
		func(sa Application) {
			wg.Go(func() {
				sa.Run(ctx)
			})
		}(sa)
	}
	if a.application != nil {
		wg.Go(func() {
			a.application.Run(ctx)
		})
	}
	wg.Wait()
}

func (a *SubApplication) Reload(ctx context.Context) { // sub first, then main
	wg := conc.NewWaitGroup()
	for _, sa := range a.subApplications {
		func(sa Application) {
			wg.Go(func() {
				sa.Reload(ctx)
			})
		}(sa)
	}
	wg.Wait()

	if a.application != nil {
		a.application.Reload(ctx)
	}
}

func (a *SubApplication) Close(ctx context.Context) { // main first, then sub
	if a.application != nil {
		a.application.Close(ctx)
	}

	wg := conc.NewWaitGroup()
	for _, sa := range a.subApplications {
		func(sa Application) {
			wg.Go(func() {
				sa.Close(ctx)
			})
		}(sa)
	}
	wg.Wait()

}
