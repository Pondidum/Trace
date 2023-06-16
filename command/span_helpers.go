package command

import (
	"context"
	"strconv"
	"time"
	"trace/tracing"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func createRootSpan(tp trace.Tracer, traceParent string, props map[string]string) (trace.Span, error) {
	nano := props["start"]
	i, err := strconv.ParseInt(nano, 10, 64)
	if err != nil {
		return nil, err
	}
	start := time.Unix(0, i)

	_, span := tp.Start(context.Background(), props["name"], trace.WithTimestamp(start))

	return span, nil
}

func createSpan(tp trace.Tracer, traceParent string, props map[string]string) (trace.Span, error) {
	nano := props["start"]
	i, err := strconv.ParseInt(nano, 10, 64)
	if err != nil {
		return nil, err
	}
	start := time.Unix(0, i)

	// cli carrier traceParent
	ctx := tracing.WithTraceParent(context.Background(), traceParent)
	_, span := tp.Start(ctx, props["name"], trace.WithTimestamp(start))

	return span, nil
}

func finishSpan(span trace.Span, finish int64) {
	span.End(trace.WithTimestamp(time.Unix(0, finish)))
}

func applyProps(span trace.Span, props map[string]string) {

	delete(props, "name")
	delete(props, "start")

	attrs := tracing.AttributesFromMap(props)

	span.SetAttributes(attrs...)
}

func applyStatus(span trace.Span, flag *pflag.Flag) {

	if flag != nil && flag.Changed {

		value := flag.Value.String()
		if value == "unset" {
			value = ""
		}

		span.SetStatus(codes.Error, value)
	} else {
		span.SetStatus(codes.Ok, "")
	}
}
