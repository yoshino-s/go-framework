package http

import (
	"context"
	"net/http"
	"sort"
	"sync/atomic"

	"github.com/go-errors/errors"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
)

var _ application.Application = (*Handler)(nil)
var _ http.Handler = (*Handler)(nil)

type Handler struct {
	*application.EmptyApplication
	*echo.Echo
	config

	Ready  *atomic.Bool
	Health *atomic.Bool
}

func New() *Handler {
	h := &Handler{
		EmptyApplication: application.NewEmptyApplication("HTTPHandler"),
		Echo:             echo.New(),
		config:           config{},
		Ready:            &atomic.Bool{},
		Health:           &atomic.Bool{},
	}

	return h
}

func (h *Handler) Configuration() configuration.Configuration {
	return &h.config
}

func (h *Handler) Setup(ctx context.Context) {
	h.EmptyApplication.Setup(ctx)

	h.Ready.Store(true)
	h.Health.Store(true)

	h.HideBanner = true
	h.HidePort = true

	h.Echo.Logger = toEchoLogger(h.EmptyApplication.Logger)

	if h.config.BehindProxy {
		h.IPExtractor = echo.ExtractIPFromXFFHeader()
		h.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if h := c.Request().Header.Get("X-Forwarded-Host"); h != "" {
					c.Request().Host = h
				}
				return next(c)
			}
		})
	}

	h.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		var message interface{}
		httpError := &echo.HTTPError{}
		if errors.As(err, &httpError) {
			code = httpError.Code
			message = httpError.Message
			if message == nil {
				message = http.StatusText(code)
			}
		} else {
			message = err.Error()
		}
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, map[string]interface{}{
					"message": message,
					"code":    code,
				})
			}
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}

	h.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if h.config.Feature.Has(FeatureVersion) && c.Request().URL.Path == "/-/version" {
				return c.String(http.StatusOK, common.Version)
			}

			if h.config.Feature.Has(FeatureHealth) && c.Request().URL.Path == "/-/healthz" {
				if h.Health.Load() {
					return c.String(http.StatusOK, "OK")
				} else {
					return c.String(http.StatusServiceUnavailable, "NG")
				}
			}

			if h.config.Feature.Has(FeatureReady) && c.Request().URL.Path == "/-/readyz" {
				if h.Ready.Load() {
					return c.String(http.StatusOK, "OK")
				}
				return c.String(http.StatusServiceUnavailable, "NG")
			}

			if h.config.Feature.Has(FeatureMetrics) && c.Request().URL.Path == "/-/metrics" {
				echoprometheus.NewHandler()(c)
			}
			return next(c)
		}
	})

	if h.config.Log {
		h.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				h.EmptyApplication.Logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
				)
				return nil
			},
		}))
	}

	if h.config.Otel {
		h.Echo.Use(otelecho.Middleware("http"))
	}

	if h.config.ResponseTraceId {
		h.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				traceID := trace.SpanFromContext(c.Request().Context()).SpanContext().TraceID()
				if traceID.IsValid() {
					c.Response().Header().Set("X-Trace-ID", traceID.String())
				}
				return next(c)
			}
		})
	}
}

func (h *Handler) Run(context.Context) {
	routers := h.Routes()
	sort.Slice(routers, func(i, j int) bool {
		if routers[i].Path == routers[j].Path {
			return routers[i].Method < routers[j].Method
		} else {
			return routers[i].Path < routers[j].Path
		}
	})
	for _, router := range routers {
		h.EmptyApplication.Logger.Sugar().Debugf("%s %s", router.Method, router.Path)
	}
	if err := h.StartH2CServer(h.config.ListenAddr, &http2.Server{}); err != nil && err != http.ErrServerClosed {
		h.EmptyApplication.Logger.Error("failed to start server", zap.Error(err))
	}
}

func (h *Handler) Close(c context.Context) {
	h.Ready.Store(false)
	h.Health.Store(false)
	if err := h.Echo.Shutdown(c); err != nil {
		h.EmptyApplication.Logger.Error("failed to close server", zap.Error(err))
	}
}

type echoContextKey struct{}

func (h *Handler) HandleGrpc(pattern string, handler http.Handler) {
	h.POST(pattern+"*", func(c echo.Context) error {
		req := c.Request().WithContext(context.WithValue(c.Request().Context(), echoContextKey{}, c))
		handler.ServeHTTP(c.Response(), req)
		return nil
	})
}

func EchoContextFromContext(ctx context.Context) (echo.Context, bool) {
	c, ok := ctx.Value(echoContextKey{}).(echo.Context)
	return c, ok
}
