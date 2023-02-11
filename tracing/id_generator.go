package tracing

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel/trace"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type fixedIdGenerator struct {
	trace trace.TraceID
	span  trace.SpanID
}

func (g *fixedIdGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	return g.trace, g.span
}

func (g *fixedIdGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	return g.span
}

func ContinueExisting(traceParent string) (tracesdk.IDGenerator, error) {
	tid, sid, err := ParseTraceParent(traceParent)
	if err != nil {
		return nil, err
	}

	return &fixedIdGenerator{tid, sid}, nil
}

// no $ at the end as a trace can have other things that we don't care about after it
var traceParentRx = regexp.MustCompile(`^[[:xdigit:]]{2}-[[:xdigit:]]{32}-[[:xdigit:]]{16}-[[:xdigit:]]{2}`)

func ParseTraceParent(traceParent string) (trace.TraceID, trace.SpanID, error) {

	if !traceParentRx.MatchString(traceParent) {
		return trace.TraceID{}, trace.SpanID{}, fmt.Errorf("invalid traceParent")
	}

	parts := strings.Split(traceParent, "-")
	if len(parts) < 3 {
		return trace.TraceID{}, trace.SpanID{}, fmt.Errorf("invalid traceParent")
	}

	tid, err := trace.TraceIDFromHex(parts[1])
	if err != nil {
		return trace.TraceID{}, trace.SpanID{}, err
	}

	sid, err := trace.SpanIDFromHex(parts[2])
	if err != nil {
		return trace.TraceID{}, trace.SpanID{}, err
	}

	return tid, sid, nil
}
