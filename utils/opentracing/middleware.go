package opentracing

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func SetTagByEcho(span opentracing.Span, c echo.Context) {
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

func SetLogByEcho(span opentracing.Span, c echo.Context) {
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
