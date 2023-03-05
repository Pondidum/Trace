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
		expectedSpanName   string
		expectedAttrs      map[string]string
	}{
		{
			description:   "no args and no environment",
			expectedExit:  1,
			expectedError: "this command takes at least 1 argument",
		},
		{
			description:  "trace parent and no command",
			args:         []string{tracing.NewTraceParent()},
			expectedExit: 1,
		},
		{
			description:  "trace parent and command without double dash",
			args:         []string{tracing.NewTraceParent(), "echo", "hello"},
			expectedExit: 0,
		},
		{
			description:    "trace parent and command with double dash",
			args:           []string{tracing.NewTraceParent(), "--", "echo", "hello", "world"},
			expectedExit:   0,
			expectedScript: `"echo" "hello" "world"`,
		},
		{
			description:  "trace parent from env and no command",
			traceParent:  tracing.NewTraceParent(),
			args:         []string{},
			expectedExit: 1,
		},
		{
			description:      "trace parent and command without double dash",
			traceParent:      tracing.NewTraceParent(),
			args:             []string{"echo", "hello", "world"},
			expectedExit:     0,
			expectedScript:   `"echo" "hello" "world"`,
			expectedSpanName: "echo hello world",
		},
		{
			description:      "trace parent and command with double dash",
			traceParent:      tracing.NewTraceParent(),
			args:             []string{"--", "echo", "hello", "world"},
			expectedExit:     0,
			expectedScript:   `"echo" "hello" "world"`,
			expectedSpanName: "echo hello world",
		},
		{
			description:  "failing command",
			traceParent:  tracing.NewTraceParent(),
			args:         []string{"--", "exit", "1"},
			expectedExit: 1,
		},
		{
			description:  "failing command with custom exit code",
			traceParent:  tracing.NewTraceParent(),
			args:         []string{"--", "exit", "8"},
			expectedExit: 8,
		},
		{
			description:      "name flag specified",
			traceParent:      tracing.NewTraceParent(),
			args:             []string{"--name", "different", "--", "echo", "hello"},
			expectedExit:     0,
			expectedSpanName: "different",
		},
		{
			description:  "attributes specified",
			traceParent:  tracing.NewTraceParent(),
			args:         []string{"--attr", "flag_enabled=true", "--", "echo", "hello"},
			expectedExit: 0,
			expectedAttrs: map[string]string{
				"flag_enabled": "true",
			},
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

			if tc.expectedSpanName != "" {
				span := exporter.Spans[0]

				assert.Equal(t, tc.expectedSpanName, span.Name())
			}

			if tc.expectedSpanStatus != codes.Unset {
				span := exporter.Spans[0]
				assert.Equal(t, tc.expectedSpanStatus, span.Status().Code)
			}

			if len(tc.expectedAttrs) > 0 {
				span := exporter.Spans[0]
				attrs := mapFromAttributes(span.Attributes())

				for key, val := range tc.expectedAttrs {
					assert.Equal(t, val, attrs[key])
				}
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
