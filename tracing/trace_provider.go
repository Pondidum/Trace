package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func CreateTraceProvider(ctx context.Context, g tracesdk.IDGenerator, exporter tracesdk.SpanExporter) (*tracesdk.TracerProvider, error) {

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSyncer(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("github actions"), // fill this from env later
			semconv.ServiceVersionKey.String("1.0.0"),       // the version of gha if availble?
			// other env attributes
		)),
		tracesdk.WithIDGenerator(g),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}
