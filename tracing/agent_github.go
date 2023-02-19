package tracing

import (
	"os"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func githubAttributes() []attribute.KeyValue {
	osType := strings.ToLower(os.Getenv("RUNNER_OS"))
	if osType == "macos" {
		osType = "darwin"
	}

	return []attribute.KeyValue{
		semconv.ServiceNameKey.String(os.Getenv("GITHUB_JOB")),
		semconv.ServiceVersionKey.String("1.0.0"),
		attribute.String("ci.provider", "github_actions"),
		fromEnv("GITHUB_ACTOR", "github.actor"),
		fromEnv("GITHUB_EVENT_NAME", "github.event"),
		fromEnv("GITHUB_REF_NAME", "github.ref.name"),
		fromEnv("GITHUB_REF_TYPE", "github.ref.type"),
		fromEnv("GITHUB_REPOSITORY", "github.repository"),
		fromEnv("GITHUB_RUN_ATTEMPT", "github.run.attempt"),
		fromEnv("GITHUB_RUN_NUMBER", "github.run.number"),
		fromEnv("GITHUB_RUN_ID", "github.run.id"),
		fromEnv("GITHUB_SHA", "github.sha"),
		semconv.OSTypeKey.String(osType),
		fromEnv("RUNNER_ARCH", "os.arch"),
		fromEnv("RUNNER_NAME", "github.runner.type"),
	}
}
