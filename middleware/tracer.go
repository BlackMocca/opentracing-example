package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

func (m *GoMiddleware) SetTracer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		spanName := fmt.Sprintf("%s %s %s", c.Scheme(), c.Request().Method, c.Path())
		span, ctx := opentracing.StartSpanFromContext(ctx, spanName)
		setTagByEcho(span, c)
		defer span.Finish()

		newReq := c.Request().WithContext(ctx)
		c.SetRequest(newReq)

		return next(c)
	}
}

func setTagByEcho(span opentracing.Span, c echo.Context) {
	var isError = false
	if c.Response().Status > http.StatusNoContent && c.Response().Status != http.StatusConflict {
		isError = true
	}

	span.SetTag("host", c.Request().Host)
	span.SetTag("User-Agent", c.Request().Header.Get("User-Agent"))
	span.SetTag("http.method", c.Request().Method)
	span.SetTag("http.status_code", c.Response().Status)
	span.SetTag("http.url", c.Path())
	span.SetTag("error", isError)
}
