package command

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"
	"trace/tracing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/codes"
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
	}, codes.Ok, ""))

	span := exporter.Spans[0]
	assert.Equal(t, trace.String(), span.SpanContext().TraceID().String())
	assert.Equal(t, trace.String(), span.Parent().TraceID().String())
	assert.Equal(t, parentSpan.String(), span.Parent().SpanID().String())

}

func TestSpanFinish(t *testing.T) {
	trace := tracing.NewTraceParent()
	traceId, parentSpan, _ := tracing.ParseTraceParent(trace)

	s := startTestSpan(trace, "--attr", "at_start=true")
	_, spanId, _ := tracing.ParseTraceParent(s)

	addAttributes(s, "cached=false", "cache_size=5")

	// finish the trace 10 seconds later
	span := finishTestSpan(s, "--attr", "at_finish=true")

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

func TestSpanErrorFlag(t *testing.T) {
	trace := tracing.NewTraceParent()

	cases := []struct {
		commandLine     []string
		expectedMessage string
		expectedStatus  codes.Code
	}{
		{
			commandLine:     []string{},
			expectedMessage: "",
			expectedStatus:  codes.Ok,
		},
		{
			commandLine:     []string{"--error"},
			expectedMessage: "",
			expectedStatus:  codes.Error,
		},
		{
			commandLine:     []string{"--error=oh no"},
			expectedMessage: "oh no",
			expectedStatus:  codes.Error,
		},
	}

	for _, tc := range cases {
		t.Run(strings.Join(tc.commandLine, "-"), func(t *testing.T) {

			s := startTestSpan(trace)
			span := finishTestSpan(s, tc.commandLine...)

			assert.Equal(t, tc.expectedStatus, span.Status().Code)
			assert.Equal(t, tc.expectedMessage, span.Status().Description)
		})
	}

}
