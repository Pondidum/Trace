package command

import (
	"strings"
	"time"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"go.opentelemetry.io/otel/sdk/trace"
)

func startTrace() string {
	return tracing.NewTraceParent()
}

func startSpan(trace string, extra ...string) string {

	ui := cli.NewMockUi()
	start, _ := NewGroupStartCommand(ui)
	start.Run(append([]string{"tests", trace}, extra...))
	tp := strings.TrimSpace(ui.OutputWriter.String())

	return tp
}

func finishSpan(span string, extra ...string) trace.ReadOnlySpan {
	ui := cli.NewMockUi()
	exporter := tracing.NewMemoryExporter()
	cmd, _ := NewGroupFinishCommand(ui)
	cmd.testSpanExporter = exporter
	cmd.now = func() int64 { return time.Now().Add(10 * time.Second).UnixNano() }

	cmd.Run(append([]string{span}, extra...))

	return exporter.Spans[0]
}

func addAttributes(span string, pairs ...string) {
	ui := cli.NewMockUi()
	cmd, _ := NewAttrCommand(ui)
	cmd.Run(append([]string{span}, pairs...))

}
