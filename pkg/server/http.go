package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HTTPConfig contains the config for an HTTP server.
type HTTPConfig struct {
	debug                bool
	host                 string
	port                 int
	readTimeout          time.Duration
	writeTimeout         time.Duration
	allowOrigins         []string
	customErrorHandler   echo.HTTPErrorHandler
	customRecoverLogFunc middleware.LogErrorFunc
}

// SetDebug set the debug.
func (h *HTTPConfig) SetDebug(d bool) *HTTPConfig {
	h.debug = d
	return h
}

// SetCustomErrorHandler set the custom error handler.
func (h *HTTPConfig) SetCustomErrorHandler(fn echo.HTTPErrorHandler) *HTTPConfig {
	h.customErrorHandler = fn
	return h
}

// SetCustomRecoverLogFunc set the custom recover log function.
func (h *HTTPConfig) SetCustomRecoverLogFunc(fn middleware.LogErrorFunc) *HTTPConfig {
	h.customRecoverLogFunc = fn
	return h
}

// SetAddress set the address.
func (h *HTTPConfig) SetAddress(host string, port int) *HTTPConfig {
	h.host = host
	h.port = port
	return h
}

// SetReadTimeout set the read timeout.
func (h *HTTPConfig) SetReadTimeout(d time.Duration) *HTTPConfig {
	h.readTimeout = d
	return h
}

// SetWriteTimeout set the write timeout.
func (h *HTTPConfig) SetWriteTimeout(d time.Duration) *HTTPConfig {
	h.writeTimeout = d
	return h
}

// C returns a default HTTP config.
func C() *HTTPConfig {
	return &HTTPConfig{
		host:         "0.0.0.0",
		port:         8080,
		readTimeout:  60 * time.Second,
		writeTimeout: 60 * time.Second,
	}
}

// HTTP contains the dependencies for an HTTP server.
type HTTP struct {
	cfg  *HTTPConfig
	Echo *echo.Echo
}

// NewHTTP creates a new HTTP server.
func NewHTTP() *HTTP {
	cfg := C()

	return NewHTTPWithConfig(cfg)
}

// NewHTTPWithConfig creates a new HTTP server with config.
func NewHTTPWithConfig(cfg *HTTPConfig) *HTTP {
	e := echo.New()

	// config basic setting.
	e.Server.Addr = fmt.Sprintf("%s:%d", cfg.host, cfg.port)
	e.Server.ReadTimeout = cfg.readTimeout
	e.Server.WriteTimeout = cfg.writeTimeout
	e.Debug = cfg.debug
	if cfg.customErrorHandler != nil {
		e.HTTPErrorHandler = cfg.customErrorHandler
	}
	e.Use(
		middleware.Logger(),
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize:       2 << 10,
			DisableStackAll: true,
			LogErrorFunc:    cfg.customRecoverLogFunc,
		}),
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				u, _ := uuid.NewV7()
				return u.String()
			},
		}),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				return strings.Contains(c.Request().URL.Path, "swagger")
			},
		}),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: cfg.allowOrigins,
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		}),
	)

	return &HTTP{
		cfg:  cfg,
		Echo: e,
	}
}

// Start starts the HTTP server.
func StartHTTP(h *HTTP) error {
	go func() {
		if err := h.Echo.StartServer(h.Echo.Server); err != nil {
			if err == http.ErrServerClosed {
				return
			}

			fmt.Printf("failed to start HTTP server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.Echo.Shutdown(ctx); err != nil {
		fmt.Printf("⇨ http server shutting down error: %v\n", err)
	}

	return nil
}

// MiddlewareSkipper return new skipper for echo router.
func MiddlewareSkipper(patterns ...string) middleware.Skipper {
	return func(c echo.Context) bool {
		for _, p := range patterns {
			if matched, _ := doublestar.Match(p, c.Path()); matched {
				return true
			}
		}
		return false
	}
}
