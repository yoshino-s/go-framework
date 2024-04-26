package common

import "context"

type ContextModifier func(ctx context.Context) (context.Context, error)
