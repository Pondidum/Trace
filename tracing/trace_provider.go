package tracing

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func CreateTraceProvider(ctx context.Context, g tracesdk.IDGenerator, exporter tracesdk.SpanExporter) (*tracesdk.TracerProvider, error) {

	attrs := detectRunner()

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSyncer(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attrs...,
		)),
		tracesdk.WithIDGenerator(g),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

func detectRunner() []attribute.KeyValue {
	if v := os.Getenv("GITHUB_ACTIONS"); v == "true" {
		return githubAttributes()
	}

	if v := os.Getenv("TEAMCITY_VERSION"); v != "" {
		return teamcityAttributes()
	}

	if v := os.Getenv("GITLAB_CI"); v == "true" {
		return gitlabAttributes()
	}

	return defaultAttributes()
}

func fromEnv(env string, name string) attribute.KeyValue {
	return attribute.String(name, os.Getenv(env))
}
