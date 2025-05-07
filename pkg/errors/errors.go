package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/cirius-go/portfolio-server/util"
)

var (
	// Alias for errors.Is
	Is = errors.Is
	// Alias for errors.As
	As = errors.As
	// Alias for errors.Unwrap
	Unwrap = errors.Unwrap
)

// ErrorType represents an error type.
// ENUM(internal,invalid_request,conflict,not_found,unauthorized,forbidden,unknown,upstream)
//
//go:generate go-enum --marshal --names --values --ptr
type ErrorType string

// AppError represents an error.
type AppError struct {
	Type     ErrorType      `json:"type"`
	Code     int            `json:"-"`
	Message  any            `json:"message"`
	Internal error          `json:"-"` // Stores the error returned by an external dependency
	meta     map[string]any `json:"-"`
}

// New creates a new app error.
func New(t ErrorType, code int, internalErr error, msg string, args ...any) *AppError {
	e := &AppError{
		Type:     t,
		Code:     code,
		Internal: internalErr,
		Message:  fmt.Sprintf(msg, args...),
	}

	return e
}

func (e *AppError) SetMeta(key string, value any) *AppError {
	newMeta := map[string]any{}
	for k, v := range e.meta {
		newMeta[k] = v
	}
	newMeta[key] = value
	return &AppError{
		Type:     e.Type,
		Code:     e.Code,
		Message:  e.Message,
		Internal: e.Internal,
		meta:     newMeta,
	}
}

// WithInternal clones the error with the internal error.
func (e *AppError) WithInternal(err error) *AppError {
	return &AppError{
		Type:     e.Type,
		Code:     e.Code,
		Message:  e.Message,
		Internal: err,
	}
}

// SetInternal sets the internal error.
// INFO: This is an alias for WithInternal.
func (e *AppError) SetInternal(err error) *AppError {
	return e.WithInternal(err)
}

// Error returns the error message.
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewInternal creates a new internal error.
func NewInternal(err error, msg string, args ...any) *AppError {
	return New(ErrorTypeInternal, http.StatusInternalServerError, err, msg, args...)
}

// NewInvalidRequest creates a new invalid request error.
func NewInvalidRequest(err error, msg string, args ...any) *AppError {
	return New(ErrorTypeInvalidRequest, http.StatusBadRequest, err, msg, args...)
}

// NewNotFound creates a new not found error.
func NewNotFound(err error, msg string, args ...any) *AppError {
	return New(ErrorTypeNotFound, http.StatusNotFound, err, msg, args...)
}

// NewConflict creates a new conflict error.
func NewConflict(err error, msg string, args ...any) *AppError {
	return New(ErrorTypeConflict, http.StatusConflict, err, msg, args...)
}

// NewUnauthorized creates a new unauthorized error.
func NewUnauthorized(err error, msg string, args ...any) *AppError {
	return New(ErrorTypeUnauthorized, http.StatusUnauthorized, err, msg, args...)
}

// NewForbidden creates a new forbidden error.
func NewForbidden(err error, msg string, args ...any) *AppError {
	return New(ErrorTypeForbidden, http.StatusForbidden, err, msg, args...)
}

// NewUpstream creates a new upstream error.
func NewUpstream(err error, msg string, args ...any) *AppError {
	return New(ErrorTypeUpstream, http.StatusBadGateway, err, msg, args...)
}

// NewUnknown creates a new unknown error.
func NewUnknown(err error, msg string, args ...any) *AppError {
	return New(ErrorTypeUnknown, http.StatusInternalServerError, err, msg, args...)
}

// FromEchoError converts an echo error to an app error.
func FromEchoError(err error) *AppError {
	if err == nil {
		return nil
	}
	if e, ok := err.(*AppError); ok {
		return e
	}

	if e, ok := err.(*echo.HTTPError); ok {
		ae := New(ErrorTypeUnknown, e.Code, e.Internal, "")
		ae.Message = e.Message
		return ae
	}
	return NewUnknown(err, "Unknown error")
}

type ErrorHandlerConfig struct {
	reqHeaderDebugFlag string
}

func NewErrorHandlerConfig() *ErrorHandlerConfig {
	return &ErrorHandlerConfig{
		reqHeaderDebugFlag: "X-Debug",
	}
}

func (e *ErrorHandlerConfig) WithReqHeaderDebugFlag(flag string) *ErrorHandlerConfig {
	e.reqHeaderDebugFlag = flag
	return e
}

func CreateEchoErrorHandler(cfgs ...*ErrorHandlerConfig) echo.HTTPErrorHandler {
	var cfg = util.IfNull(NewErrorHandlerConfig(), cfgs...)
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var (
			srvDebug   = c.Echo().Debug
			reqDebug   = util.StrBool(c.Request().Header.Get(cfg.reqHeaderDebugFlag))
			debug      = util.IfZero(srvDebug, reqDebug)
			statusCode = http.StatusInternalServerError
			res        map[string]any
			meta       map[string]any
		)

		type Validation struct {
			Key      string `json:"key,omitempty"`
			Field    string `json:"field,omitempty"`
			FailedOn string `json:"failed_on,omitempty"`
			Value    any    `json:"value,omitempty"`
			Param    any    `json:"param,omitempty"`
			Msg      string `json:"msg,omitempty"`
		}

		parseValidationErr := func(meta map[string]any, ierr error) {
			switch ierr := ierr.(type) {
			case validator.ValidationErrors:
				if meta == nil {
					meta = map[string]any{}
				}
				fieldErrs := make([]*Validation, 0, len(ierr))
				for _, f := range ierr {
					fieldErrs = append(fieldErrs, &Validation{
						Key:      f.StructNamespace(),
						Field:    f.Namespace(),
						FailedOn: f.Tag(),
						Value:    f.Value(),
						Param:    f.Param(),
						Msg:      f.Error(),
					})
				}
				meta["validation"] = fieldErrs
				res["meta"] = meta
			}
		}

		// underline error type
		switch uerr := err.(type) {
		case *AppError:
			meta = uerr.meta
			res = map[string]any{
				"type":    uerr.Type,
				"code":    uerr.Code,
				"message": uerr.Message,
				"meta":    meta,
			}
			statusCode = uerr.Code
			if uerr.Internal != nil {
				if debug {
					res["internal"] = uerr.Internal.Error()
				}
				parseValidationErr(meta, uerr.Internal)

			}
		case *echo.HTTPError:
			meta = map[string]any{}
			res = map[string]any{
				"type":    ErrorTypeUnknown,
				"code":    uerr.Code,
				"message": uerr.Message,
			}
			statusCode = uerr.Code
			if uerr.Internal != nil {
				if debug {
					res["internal"] = uerr.Internal.Error()
				}
				parseValidationErr(meta, uerr.Internal)
			}
		default:
			msg := "Unknown error returned without any message"
			if uerr != nil {
				msg = uerr.Error()
			}
			res = map[string]any{
				"type":    ErrorTypeUnknown,
				"code":    http.StatusInternalServerError,
				"message": msg,
			}
			statusCode = http.StatusInternalServerError
		}

		if marshalErr := c.JSON(statusCode, res); marshalErr != nil {
			fmt.Println("cannot marshal error response", marshalErr)
		}
	}
}
