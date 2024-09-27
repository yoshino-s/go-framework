package telemetry

import (
	"github.com/getsentry/sentry-go"
	"github.com/yoshino-s/go-framework/common"
	"go.uber.org/zap"
)

func (c *telemetryConfiguration) initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:                c.SentryDSN,
		Debug:              common.IsDev(),
		EnableTracing:      true,
		AttachStacktrace:   true,
		TracesSampleRate:   c.TracesSampleRate,
		ProfilesSampleRate: c.ProfilesSampleRate,
		SendDefaultPII:     true,
		Release:            common.Version,
	})
	if err != nil {
		zap.L().Error("sentry.Init failed", zap.Error(err))
	}
}

func IsSentryInitialized() bool {
	return sentry.CurrentHub().Client() != nil
}
