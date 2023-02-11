package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestIdGeneration(t *testing.T) {
	t.Parallel()

	t.Run("trace_ids", func(t *testing.T) {
		seen := map[trace.TraceID]bool{}

		for i := 0; i < 10; i++ {
			seen[NewTraceID()] = true
		}

		assert.Len(t, seen, 10)
	})

	t.Run("span_ids", func(t *testing.T) {
		seen := map[trace.SpanID]bool{}

		for i := 0; i < 10; i++ {
			seen[NewSpanID()] = true
		}

		assert.Len(t, seen, 10)
	})

}
