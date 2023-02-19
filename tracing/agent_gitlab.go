package tracing

import (
	"os"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func gitlabAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ServiceNameKey.String(os.Getenv("CI_JOB_NAME")),
		semconv.ServiceVersionKey.String(os.Getenv("CI_SERVER_VERSION")),
		attribute.String("ci.provider", "gitlab_ci"),
		fromEnv("CI_PIPELINE_SOURCE", "gitlab.pipeline.source"),
		fromEnv("CI_COMMIT_REF_NAME", "gitlab.ref.name"),
		fromEnv("CI_PROJECT_NAME", "gitlab.project"),
		fromEnv("CI_JOB_ID", "gitlab.job.id"),
		fromEnv("CI_PIPELINE_ID", "gitlab.pipeline.id"),
		fromEnv("CI_COMMIT_SHA", "gitlab.sha"),
	}
}
