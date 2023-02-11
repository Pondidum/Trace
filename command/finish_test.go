package command

import (
	"os"
	"strings"
	"testing"
	"time"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestFinishingTrace(t *testing.T) {
	t.Parallel()

	t.Run("when there is no traceparent arg or env", func(t *testing.T) {
		os.Unsetenv(TraceParentEnvVar)

		ui := cli.NewMockUi()
		cmd, _ := NewFinishCommand(ui)

		assert.Equal(t, 1, cmd.Run([]string{}))
		assert.Contains(t, ui.ErrorWriter.String(), "this command takes 1 argument: traceparent")
	})

	t.Run("when there is a traceparent arg", func(t *testing.T) {
		os.Unsetenv(TraceParentEnvVar)
		tp := startTestTrace()

		cmd, ui, _ := createTestFinishCommand()

		assert.Equal(t, 0, cmd.Run([]string{tp}), ui.ErrorWriter.String())
	})

	t.Run("when there is a traceparent envvar", func(t *testing.T) {
		tp := startTestTrace()
		os.Setenv(TraceParentEnvVar, tp)

		cmd, ui, _ := createTestFinishCommand()

		assert.Equal(t, 0, cmd.Run([]string{}), ui.ErrorWriter.String())
	})

	t.Run("all attributes are recorded", func(t *testing.T) {
		startTime := time.Now().UnixNano()
		endTime := time.Now().Add(10 * time.Second).UnixNano()

		// start a trace
		ui := cli.NewMockUi()
		start, _ := NewStartCommand(ui)
		start.now = func() int64 { return startTime }
		start.Run([]string{"tests"})
		tp := strings.TrimSpace(ui.OutputWriter.String())

		// finish the trace 10 seconds later
		cmd, _, exporter := createTestFinishCommand()
		cmd.now = func() int64 { return endTime }
		assert.Equal(t, 0, cmd.Run([]string{tp}), ui.ErrorWriter.String())

		span := exporter.Spans[0]
		assert.Len(t, exporter.Spans, 1)
		assert.Equal(t, "tests", span.Name())
		assert.Equal(t, "trace-cli", span.InstrumentationLibrary().Name)
		assert.Equal(t, startTime, span.StartTime().UnixNano())
		assert.Equal(t, endTime, span.EndTime().UnixNano())
	})

}

func createTestFinishCommand() (*FinishCommand, *cli.MockUi, *tracing.MemoryExporter) {

	ui := cli.NewMockUi()
	exporter := tracing.NewMemoryExporter()

	cmd, _ := NewFinishCommand(ui)
	cmd.testTracerProvider = tracesdk.NewTracerProvider(
		tracesdk.WithSyncer(exporter),
	)

	return cmd, ui, exporter
}

func startTestTrace() string {
	ui := cli.NewMockUi()
	cmd, _ := NewStartCommand(ui)
	cmd.Run([]string{"tests"})

	return strings.TrimSpace(ui.OutputWriter.String())
}
