package http

import (
	"context"
	"net/http"
	"sort"
	"sync/atomic"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/telemetry"
	"go.uber.org/zap"
)

var _ application.Application = (*Handler)(nil)
var _ http.Handler = (*Handler)(nil)

type Handler struct {
	*echo.Echo
	config Config
	logger *zap.Logger
	Ready  *atomic.Bool
	Health *atomic.Bool
}

var SkipLog = map[string]bool{
	"/healthz": true,
	"/readyz":  true,
	"/version": true,
}

func New(config Config) (*Handler, error) {
	e := echo.New()
	logger := zap.L()

	h := &Handler{
		Echo:   e,
		config: config,
		Ready:  &atomic.Bool{},
		Health: &atomic.Bool{},
		logger: logger,
	}

	if config.Debug {
		e.Debug = true
		logger.Info("debug mode enabled")
	}

	if telemetry.IsSentryInitialized() {
		e.Use(sentryecho.New(sentryecho.Options{}))
	}

	if config.Log {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if _, ok := SkipLog[v.URI]; ok {
					return nil
				}
				logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
				)
				return nil
			},
		}))
	}

	if config.Feature.Has(FeatureVersion) {
		e.GET("/version", func(c echo.Context) error {
			return c.String(http.StatusOK, "1")
		})
	}

	if config.Feature.Has(FeatureHealth) {
		e.GET("/healthz", func(c echo.Context) error {
			if h.Health.Load() {
				return c.String(http.StatusOK, "OK")
			} else {
				return c.String(http.StatusServiceUnavailable, "NG")
			}
		})
	}

	if config.Feature.Has(FeatureReady) {
		e.GET("/readyz", func(c echo.Context) error {
			if h.Ready.Load() {
				return c.String(http.StatusOK, "OK")
			}
			return c.String(http.StatusServiceUnavailable, "NG")
		})
	}

	if config.Feature.Has(FeatureMetrics) {
		e.GET("/metrics", echoprometheus.NewHandler())
	}

	return h, nil
}

func (h *Handler) Setup(context.Context) {
	h.Ready.Store(true)
	h.Health.Store(true)
}

func (h *Handler) Reload(context.Context) {
	h.Health.Store(false)
	defer h.Health.Store(true)
}

func (h *Handler) Run(context.Context) {
	if h.config.Debug {
		routers := h.Routes()
		sort.Slice(routers, func(i, j int) bool {
			if routers[i].Path == routers[j].Path {
				return routers[i].Method < routers[j].Method
			} else {
				return routers[i].Path < routers[j].Path
			}
		})
		for _, router := range routers {
			h.logger.Sugar().Debugf("%s %s", router.Method, router.Path)
		}
	}
	if err := h.Start(h.config.ListenAddr); err != nil {
		h.logger.Error("failed to start server", zap.Error(err))
	}
}

func (h *Handler) Close(context.Context) {
	h.Ready.Store(false)
	h.Health.Store(false)
	if err := h.Echo.Close(); err != nil {
		h.logger.Error("failed to close server", zap.Error(err))
	}
}
