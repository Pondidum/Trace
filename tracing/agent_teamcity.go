package tracing

import (
	"os"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func teamcityAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ServiceNameKey.String(os.Getenv("TEAMCITY_BUILDCONF_NAME")),
		semconv.ServiceVersionKey.String(os.Getenv("TEAMCITY_VERSION")),
		attribute.String("ci.provider", "teamcity"),
		fromEnv("TEAMCITY_PROJECT_NAME", "teamcity.project.name"),
		fromEnv("BUILD_NUMBER", "teamcity.run.number"),
	}
}
