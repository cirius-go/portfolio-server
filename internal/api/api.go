package api

import (
	"context"
	"net/http"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/cirius-go/portfolio-server/util"
)

// API represents the base API.
type API struct{}

// APIOptions contains the options for the handler function.
type APIOptions struct {
	successStatusCode int
}

// O creates a new API options.
func O() *APIOptions {
	return &APIOptions{}
}

// SuccessStatus sets the custom status code for the success response.
func (o *APIOptions) SuccessStatus(code int) *APIOptions {
	o.successStatusCode = code
	return o
}

// ServiceHandlerFunc represents the service handler function with request and response models.
type ServiceHandlerFunc[Rq, Rp any] func(context.Context, *Rq) (*Rp, error)

// ServiceNoReqHandlerFunc represents the service handler function with response model only.
type ServiceNoReqHandlerFunc[Rp any] func(context.Context) (*Rp, error)

// ServiceNoResHandlerFunc represents the service handler function with request model only.
type ServiceNoResHandlerFunc[Rq any] func(context.Context, *Rq) error

// MakeJSONHandler trigger resolver with echo context.
func MakeJSONHandler[Rq, Rp any](c echo.Context, fn ServiceHandlerFunc[Rq, Rp], opts ...*APIOptions) error {
	return JSONHandlerFunc(fn)(c)
}

// JSONHandlerFunc wrap the service method with json parsers.
func JSONHandlerFunc[Rq, Rp any](fn ServiceHandlerFunc[Rq, Rp], opts ...*APIOptions) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			opt *APIOptions

			rq  = new(Rq)
			ctx = c.Request().Context()

			err error
		)

		if len(opts) > 0 {
			opt = opts[0]
		} else {
			opt = O()
		}

		if err := c.Bind(rq); err != nil {
			return err
		}

		res, err := fn(ctx, rq)
		if err != nil {
			return err
		}

		successStatus := util.IfZero(http.StatusOK, opt.successStatusCode)
		return c.JSON(successStatus, res)
	}
}

// NoReqHandlerFunc wrap the service method with json parsers.
func NoReqHandlerFunc[Rp any](fn ServiceNoReqHandlerFunc[Rp], opts ...*APIOptions) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			opt *APIOptions

			ctx = c.Request().Context()

			err error
		)

		if len(opts) > 0 {
			opt = opts[0]
		} else {
			opt = O()
		}

		res, err := fn(ctx)
		if err != nil {
			return err
		}

		successStatus := util.IfZero(http.StatusOK, opt.successStatusCode)
		return c.JSON(successStatus, res)
	}
}

// NoResHandlerFunc wrap the service method with json parsers.
func NoResHandlerFunc[Rq any](fn ServiceNoResHandlerFunc[Rq], opts ...*APIOptions) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			opt *APIOptions

			rq  = new(Rq)
			ctx = c.Request().Context()

			err error
		)

		if len(opts) > 0 {
			opt = opts[0]
		} else {
			opt = O()
		}

		if err := c.Bind(rq); err != nil {
			return err
		}

		if err = fn(ctx, rq); err != nil {
			return err
		}

		successStatus := util.IfZero(http.StatusOK, opt.successStatusCode)
		return c.NoContent(successStatus)
	}
}

// PathSkipper return new skipper for echo router.
func PathSkipper(patterns ...string) middleware.Skipper {
	return func(c echo.Context) bool {
		for _, p := range patterns {
			if matched, _ := doublestar.Match(p, c.Path()); matched {
				return true
			}
		}
		return false
	}
}
