package tracing

import (
	"context"

	"go.opentelemetry.io/otel/sdk/trace"
)

func NewMemoryExporter() *MemoryExporter {
	return &MemoryExporter{}
}

type MemoryExporter struct {
	Spans []trace.ReadOnlySpan
}

func (e *MemoryExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	e.Spans = append(e.Spans, spans...)

	return nil
}

func (e *MemoryExporter) Shutdown(ctx context.Context) error {
	return nil
}
