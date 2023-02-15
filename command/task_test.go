package command

import (
	"os"
	"testing"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TestTaskArgumentParsing(t *testing.T) {

	cases := []struct {
		description        string
		traceParent        string
		args               []string
		expectedExit       int
		expectedError      string
		expectedScript     string
		expectedSpanStatus codes.Code
	}{
		{
			description:   "no args and no environment",
			expectedExit:  1,
			expectedError: "this command takes at least 1 argument",
		},
		{
			description:  "trace parent and no command",
			args:         []string{NewTraceParent()},
			expectedExit: 1,
		},
		{
			description:  "trace parent and command without double dash",
			args:         []string{NewTraceParent(), "echo", "hello"},
			expectedExit: 0,
		},
		{
			description:    "trace parent and command with double dash",
			args:           []string{NewTraceParent(), "--", "echo", "hello", "world"},
			expectedExit:   0,
			expectedScript: `"echo" "hello" "world"`,
		},
		{
			description:  "trace parent from env and no command",
			traceParent:  NewTraceParent(),
			args:         []string{},
			expectedExit: 1,
		},
		{
			description:    "trace parent and command without double dash",
			traceParent:    NewTraceParent(),
			args:           []string{"echo", "hello", "world"},
			expectedExit:   0,
			expectedScript: `"echo" "hello" "world"`,
		},
		{
			description:    "trace parent and command with double dash",
			traceParent:    NewTraceParent(),
			args:           []string{"--", "echo", "hello", "world"},
			expectedExit:   0,
			expectedScript: `"echo" "hello" "world"`,
		},
		{
			description:  "failing command",
			traceParent:  NewTraceParent(),
			args:         []string{"--", "exit", "1"},
			expectedExit: 1,
		},
		{
			description:  "failing command with custom exit code",
			traceParent:  NewTraceParent(),
			args:         []string{"--", "exit", "8"},
			expectedExit: 8,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			os.Setenv(TraceParentEnvVar, tc.traceParent)
			if tc.args == nil {
				tc.args = []string{}
			}

			ui := cli.NewMockUi()
			exporter := tracing.NewMemoryExporter()
			cmd, _ := NewTaskCommand(ui)
			cmd.testSpanExporter = exporter

			exitCode := cmd.Run(tc.args)
			assert.Equal(t, tc.expectedExit, exitCode, ui.ErrorWriter.String())

			if tc.expectedError != "" {
				assert.Contains(t, ui.ErrorWriter.String(), tc.expectedError)
			}

			if tc.expectedScript != "" {
				span := exporter.Spans[0]
				attrs := mapFromAttributes(span.Attributes())

				assert.Equal(t, tc.expectedScript, attrs["shell_script"])
			}

			if tc.expectedSpanStatus != codes.Unset {
				span := exporter.Spans[0]
				assert.Equal(t, tc.expectedSpanStatus, span.Status().Code)
			}
		})
	}
}

func mapFromAttributes(attrs []attribute.KeyValue) map[string]string {
	m := make(map[string]string, len(attrs))
	for _, at := range attrs {
		m[string(at.Key)] = at.Value.AsString()
	}
	return m
}
