// Package httpres
//
//nolint:staticcheck
package httpres

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	CodeServiceNotAvailable = "ServiceNotAvailable"
	CodeDenied              = "denied"
	CodeBadRequest          = "badRequest"
)

type HTTPError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type appError struct {
	err          error
	businessCode string
	metricLabel  string
	httpStatus   int
}

func (e *appError) Error() string        { return e.err.Error() }
func (e *appError) Unwrap() error        { return e.err }
func (e *appError) HTTPStatus() int      { return e.httpStatus }
func (e *appError) BusinessCode() string { return e.businessCode }
func (e *appError) MetricLabel() string  { return e.metricLabel }

func NewHTTPError(err error, status int, businessCode, label string) error {
	return &appError{
		err:          err,
		httpStatus:   status,
		businessCode: businessCode,
		metricLabel:  label,
	}
}

func NewError(c *gin.Context, status int, code string, err error) {
	c.JSON(status, HTTPError{
		ErrorCode:    code,
		ErrorMessage: err.Error(),
	})
	c.Abort()
}

func HandleDomainError(ctx *gin.Context, err error, metric *prometheus.CounterVec) {
	var mapped interface {
		error
		HTTPStatus() int
		BusinessCode() string
		MetricLabel() string
	}

	if errors.As(err, &mapped) {
		if metric != nil {
			metric.WithLabelValues(mapped.MetricLabel()).Inc()
		}

		NewError(ctx, mapped.HTTPStatus(), mapped.BusinessCode(), err)
		return
	}

	if metric != nil {
		metric.WithLabelValues("error").Inc()
	}
	_ = ctx.Error(err).SetType(gin.ErrorTypePrivate)
	cerr := errors.New("Что-то пошло не так")
	NewError(ctx, http.StatusInternalServerError, CodeServiceNotAvailable, cerr)
}
