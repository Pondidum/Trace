package command

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"

	"go.opentelemetry.io/otel/trace"
)

var randSource *rand.Rand

// invoked by go runtime
func init() {
	if randSource != nil {
		return
	}

	var rngSeed int64
	binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	randSource = rand.New(rand.NewSource(rngSeed))
}

func NewTraceID() trace.TraceID {
	tid := trace.TraceID{}
	randSource.Read(tid[:])
	return tid
}

func NewSpanID() trace.SpanID {
	sid := trace.SpanID{}
	randSource.Read(sid[:])
	return sid
}

func NewTraceParent() string {
	return fmt.Sprintf("00-%s-%s-01", NewTraceID(), NewSpanID())
}

func AsTraceParent(tid trace.TraceID, sid trace.SpanID) string {
	return fmt.Sprintf("00-%s-%s-01", tid, sid)
}
