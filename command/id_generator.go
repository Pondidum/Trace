package command

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

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
