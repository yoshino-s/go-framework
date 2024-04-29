package telemetry

import (
	"github.com/getsentry/sentry-go"
	"github.com/yoshino-s/go-framework/common"
	"go.uber.org/zap"
)

func initSentry(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            common.IsDev(),
		EnableTracing:    true,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   true,
		Release:          common.Version,
	})
	if err != nil {
		zap.L().Error("sentry.Init failed", zap.Error(err))
	}
}

func IsSentryInitialized() bool {
	return sentry.CurrentHub().Client() != nil
}
