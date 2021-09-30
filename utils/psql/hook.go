package psql

import (
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/spf13/cast"
)

type TracingHook struct {
	tracer opentracing.Tracer
}

func NewTracingHook(tracing opentracing.Tracer) *TracingHook {
	return &TracingHook{
		tracer: tracing,
	}
}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *TracingHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if ctx != nil {
		span := opentracing.SpanFromContext(ctx)
		span, ctx = opentracing.StartSpanFromContext(ctx, "database", opentracing.ChildOf(span.Context()))
		span.LogFields(
			otlog.String("statement", query),
		)

		if args != nil && len(args) > 0 {
			var argsString = []string{}
			for index, arg := range args {
				argsString = append(argsString, fmt.Sprintf(`$$%s:%s`, cast.ToString(index+1), cast.ToString(arg)))
			}
			span.LogFields(
				otlog.String("args", strings.Join(argsString, ",")),
			)
		}
	}
	return ctx, nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *TracingHook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if ctx != nil {
		span := opentracing.SpanFromContext(ctx)
		defer span.Finish()

		span.SetTag("error", false)
	}
	return ctx, nil
}

// Hook OnError
func (h *TracingHook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	if ctx != nil {
		span := opentracing.SpanFromContext(ctx)
		defer span.Finish()

		span.SetTag("error", true)
		span.LogFields(
			otlog.Message(err.Error()),
		)
	}

	return err
}
