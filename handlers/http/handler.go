package http

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"sync/atomic"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-errors/errors"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yoshino-s/go-framework/application"
	"github.com/yoshino-s/go-framework/authentication/oidc"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
	framework_errors "github.com/yoshino-s/go-framework/errors"
	"github.com/yoshino-s/go-framework/telemetry"
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

	oidcAuthenticationRegisterFunc oidc.RegisterFunc
}

func New() *Handler {
	h := &Handler{
		EmptyApplication: application.NewEmptyApplication(),
		Echo:             echo.New(),
		config:           config{},
		Ready:            &atomic.Bool{},
		Health:           &atomic.Bool{},
	}

	return h
}

func (h *Handler) SetOIDCAuthentication(auth *oidc.OIDCAuthentication, RedirectPath string, CallbackPath string, PostProcess oidc.PostProcessFunc) error {
	oidcAuthenticationRegisterFunc, err := auth.Register(oidc.MiddlewareConfig{
		ExternalURL:  h.config.ExternalURL,
		RedirectPath: RedirectPath,
		CallbackPath: CallbackPath,
		PostProcess:  PostProcess,
	})
	if err != nil {
		return err
	}
	h.oidcAuthenticationRegisterFunc = oidcAuthenticationRegisterFunc
	return nil
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
		appError := &framework_errors.AppError{}
		if errors.As(err, &httpError) {
			code = httpError.Code
			message = httpError.Message
			if message == nil {
				message = http.StatusText(code)
			}
		} else if errors.As(err, &appError) {
			code = appError.Code()
			message = appError.Error()
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

	if telemetry.IsSentryInitialized() {
		h.Use(sentryecho.New(sentryecho.Options{}))
	}

	if h.config.Log {
		h.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if strings.HasPrefix(v.URI, "/-/") {
					return nil
				}
				h.EmptyApplication.Logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
				)
				return nil
			},
		}))
	}

	if h.config.Feature.Has(FeatureVersion) {
		h.GET("/-/version", func(c echo.Context) error {
			return c.String(http.StatusOK, common.Version)
		})
	}

	if h.config.Feature.Has(FeatureHealth) {
		h.GET("/-/healthz", func(c echo.Context) error {
			if h.Health.Load() {
				return c.String(http.StatusOK, "OK")
			} else {
				return c.String(http.StatusServiceUnavailable, "NG")
			}
		})
	}

	if h.config.Feature.Has(FeatureReady) {
		h.GET("/-/readyz", func(c echo.Context) error {
			if h.Ready.Load() {
				return c.String(http.StatusOK, "OK")
			}
			return c.String(http.StatusServiceUnavailable, "NG")
		})
	}

	if h.config.Feature.Has(FeatureMetrics) {
		h.GET("/-/metrics", echoprometheus.NewHandler())
	}

	if h.oidcAuthenticationRegisterFunc != nil {
		h.oidcAuthenticationRegisterFunc(h.Echo)
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
