package tracing

import (
	"os"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func defaultAttributes() []attribute.KeyValue {

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = os.Getenv("SHELL")
	}

	return []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String("1.0.0"),
	}
}
