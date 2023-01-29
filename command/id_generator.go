package command

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"

	"go.opentelemetry.io/otel/trace"
)

var randSource *rand.Rand

func init() {
	if randSource != nil {
		return
	}

	var rngSeed int64
	binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	randSource = rand.New(rand.NewSource(rngSeed))
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	sid := trace.SpanID{}
	_, _ = randSource.Read(sid[:])
	return sid
}

// NewIDs returns a non-zero trace ID and a non-zero span ID from a
// randomly-chosen sequence.
func NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	tid := trace.TraceID{}
	_, _ = randSource.Read(tid[:])
	sid := trace.SpanID{}
	_, _ = randSource.Read(sid[:])
	return tid, sid
}

func NewTraceParent(ctx context.Context) string {
	trace, span := NewIDs(ctx)
	return fmt.Sprintf("00-%s-%s-01", trace, span)
}
