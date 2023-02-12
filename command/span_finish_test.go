package command

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"
	"trace/tracing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestParentingSpans(t *testing.T) {

	tp := NewTraceParent()
	trace, parentSpan, _ := tracing.ParseTraceParent(tp)
	exporter := tracing.NewMemoryExporter()

	g, _ := tracing.ContinueExisting(tp)
	provider, err := tracing.CreateTraceProvider(context.Background(), g, exporter)
	assert.NoError(t, err)
	tracer := provider.Tracer("tests")

	start := time.Now()
	finish := start.Add(10 * time.Second)

	assert.NoError(t, createSpan(tracer, tp, finish.UnixNano(), map[string]string{
		"start": strconv.FormatInt(start.UnixNano(), 10),
	}))

	span := exporter.Spans[0]
	assert.Equal(t, trace.String(), span.SpanContext().TraceID().String())
	assert.Equal(t, trace.String(), span.Parent().TraceID().String())
	assert.Equal(t, parentSpan.String(), span.Parent().SpanID().String())

}

func TestSpanFinish(t *testing.T) {
	trace := startTrace()
	traceId, parentSpan, _ := tracing.ParseTraceParent(trace)

	s := startSpan(trace)
	_, spanId, _ := tracing.ParseTraceParent(s)

	// finish the trace 10 seconds later
	span := finishSpan(s)

	// debug information
	t.Log("trace :", traceId.String())
	t.Log("parent:", parentSpan.String())
	t.Log("span  :", spanId.String())

	assert.Equal(t, "tests", span.Name())
	assert.Equal(t, "trace-cli", span.InstrumentationLibrary().Name)
	assert.Equal(t, traceId.String(), span.SpanContext().TraceID().String(), "wrong trace")
	assert.Equal(t, traceId.String(), span.Parent().TraceID().String(), "wrong parent trace")
	assert.Equal(t, spanId.String(), span.SpanContext().SpanID().String(), "wrong id")
	assert.Equal(t, parentSpan.String(), span.Parent().SpanID().String(), "wrong parent id")
}

func startTrace() string {
	return NewTraceParent()
}

func startSpan(trace string) string {

	ui := cli.NewMockUi()
	start, _ := NewSpanStartCommand(ui)
	start.Run([]string{"tests", trace})
	tp := strings.TrimSpace(ui.OutputWriter.String())

	return tp
}

func finishSpan(span string) trace.ReadOnlySpan {
	ui := cli.NewMockUi()
	exporter := tracing.NewMemoryExporter()
	cmd, _ := NewSpanFinishCommand(ui)
	cmd.testSpanExporter = exporter
	cmd.now = func() int64 { return time.Now().Add(10 * time.Second).UnixNano() }

	cmd.Run([]string{span})

	return exporter.Spans[0]
}
