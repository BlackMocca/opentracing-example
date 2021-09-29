package main

import (
	"context"
	"fmt"
	"io"

	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func main() {
	/* init tracing*/

	tracer, closer := tracing.Init("hello-world")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	/* start span */
	var userId, _ = uuid.NewV4()

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "main-func")
	defer span.Finish()

	ctx = spanA(ctx, &userId)
	spanB(ctx, &userId)
}

func spanA(ctx context.Context, userId *uuid.UUID) context.Context {
	span, ctx := opentracing.StartSpanFromContext(ctx, "spanA-func")
	defer span.Finish()

	span.LogFields(
		log.Event("fetch user usecase"),
		log.String("user_id", userId.String()),
	)

	return ctx
}

func spanB(ctx context.Context, userId *uuid.UUID) context.Context {
	span, ctx := opentracing.StartSpanFromContext(ctx, "spanB-func")
	defer span.Finish()

	span.LogFields(
		log.Event("fetch user repository"),
		log.String("user_id", userId.String()),
	)

	return ctx
}
