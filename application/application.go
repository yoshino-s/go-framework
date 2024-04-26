package application

import "context"

type Application interface {
	Setup(context.Context)
	Run(context.Context)
	Reload(context.Context)
	Close(context.Context)
}
