package telemetry

func Init(config Config) {
	if config.SentryDSN != "" {
		InitSentry(config.SentryDSN)
	}
}
