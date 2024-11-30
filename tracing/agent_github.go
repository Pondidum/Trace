package tracing

import (
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

func githubAttributes() []attribute.KeyValue {
	osType := strings.ToLower(os.Getenv("RUNNER_OS"))
	if osType == "macos" {
		osType = "darwin"
	}

	serverUrl := os.Getenv("GITHUB_SERVER_URL")
	repo := os.Getenv("GITHUB_REPOSITORY")
	runId := os.Getenv("GITHUB_RUN_ID")
	attempt := os.Getenv("GITHUB_RUN_ATTEMPT")

	url := fmt.Sprintf("%s/%s/actions/runs/%s/attempts/%s", serverUrl, repo, runId, attempt)

	refName := os.Getenv("GITHUB_REF_NAME")
	refType := os.Getenv("GITHUB_REF_TYPE")
	sha := os.Getenv("GITHUB_SHA")

	return []attribute.KeyValue{
		semconv.ServiceName(os.Getenv("GITHUB_JOB")),
		semconv.ServiceVersion("1.0.0"),
		attribute.String("ci.provider", "github_actions"),
		attribute.String("github.ref.name", refName),
		attribute.String("github.ref.type", refType),
		attribute.String("github.repository", repo),
		attribute.String("github.run.attempt", attempt),
		attribute.String("github.run.id", runId),
		attribute.String("github.sha", sha),
		semconv.OSTypeKey.String(osType),
		fromEnv("GITHUB_ACTOR", "github.actor"),
		fromEnv("GITHUB_EVENT_NAME", "github.event"),
		fromEnv("GITHUB_RUN_NUMBER", "github.run.number"),
		fromEnv("RUNNER_ARCH", "os.arch"),
		fromEnv("RUNNER_NAME", "github.runner.type"),

		// cicd semconv
		semconv.CICDPipelineRunID(runId),
		semconv.CICDPipelineName(os.Getenv("GITHUB_WORKFLOW")),
		semconv.CICDPipelineTaskRunID(os.Getenv("GITHUB_ACTION")),
		// semconv.CICDPipelineTaskName(os.Getenv("GITHUB_ACTION")),
		semconv.CICDPipelineTaskRunURLFull(url),

		// vsc semconv
		semconv.VCSRepositoryRefName(refName),
		semconv.VCSRepositoryRefRevision(sha),
		attribute.String(string(semconv.VCSRepositoryRefTypeKey), refType),
		semconv.VCSRepositoryURLFull(fmt.Sprintf("%s/%s", serverUrl, repo)),
	}
}
