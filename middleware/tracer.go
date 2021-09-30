package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func (m *GoMiddleware) SetTracer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var span opentracing.Span
		var ctx = c.Request().Context()
		var spanName = fmt.Sprintf("%s %s %s", c.Scheme(), c.Request().Method, c.Path())

		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request().Header))
		switch err {
		case nil:
			/* has parent span context */
			span = opentracing.StartSpan(spanName, ext.RPCServerOption(spanCtx))
		case opentracing.ErrSpanContextNotFound:
			/* new span context */
			span, ctx = opentracing.StartSpanFromContext(ctx, spanName)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		setTagByEcho(span, c)
		setLogByEcho(span, c)
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

func setLogByEcho(span opentracing.Span, c echo.Context) {
	var paramNameM = map[string]string{}
	var paramsName = c.ParamNames()
	var paramsValue = c.ParamValues()
	var paramLogString string

	if paramsName != nil && len(paramsName) > 0 {
		for index, paramName := range paramsName {
			paramNameM[paramName] = ""
			if index <= len(paramsValue)-1 {
				paramNameM[paramName] = paramsValue[index]
			}
		}
	}

	if len(paramNameM) > 0 {
		var logs = []string{}
		for k, v := range paramNameM {
			logs = append(logs, fmt.Sprintf("%s:%s", k, v))
		}
		paramLogString = strings.Join(logs, ",")
	}

	span.LogFields(
		log.String("querystring", c.QueryString()),
		log.String("param", paramLogString),
	)
}
