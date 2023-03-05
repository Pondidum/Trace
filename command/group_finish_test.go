package command

import (
	"context"
	"strconv"
	"testing"
	"time"
	"trace/tracing"

	"github.com/stretchr/testify/assert"
)

func TestParentingSpans(t *testing.T) {

	tp := tracing.NewTraceParent()
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

	s := startSpan(trace, "--attr", "at_start=true")
	_, spanId, _ := tracing.ParseTraceParent(s)

	addAttributes(s, "cached=false", "cache_size=5")

	// finish the trace 10 seconds later
	span := finishSpan(s, "--attr", "at_finish=true")

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

	attrs := mapFromAttributes(span.Attributes())
	assert.Contains(t, attrs, "cached")
	assert.Contains(t, attrs, "cache_size")
	assert.Equal(t, "false", attrs["cached"])
	assert.Equal(t, "5", attrs["cache_size"])
	assert.Equal(t, "true", attrs["at_start"])
	assert.Equal(t, "true", attrs["at_finish"])
}
