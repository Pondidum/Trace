package tracing

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func CreateExporter(ctx context.Context) (tracesdk.SpanExporter, error) {
	exporter := os.Getenv("OTEL_TRACE_EXPORTER")

	switch exporter {
	case "none":
		return nil, nil

	case "stdout":
		return stdouttrace.New(stdouttrace.WithPrettyPrint())

	case "stderr":
		return stdouttrace.New(stdouttrace.WithPrettyPrint(), stdouttrace.WithWriter(os.Stderr))

	default:
		return otlptracegrpc.New(ctx)
	}
}
