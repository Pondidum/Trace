package tracing

import (
	"os"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func defaultAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ServiceNameKey.String(os.Getenv("SHELL")),
		semconv.ServiceVersionKey.String("1.0.0"),
	}
}
