package tracing

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestEnvironmentVariableDetection(t *testing.T) {

	os.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.namespace=some-value")

	exporter := NewMemoryExporter()
	p, err := CreateTraceProvider(context.Background(), &fixedIdGenerator{}, exporter)
	assert.NoError(t, err)

	tr := p.Tracer("test")
	_, span := tr.Start(context.Background(), "test-span")
	span.End()

	sent := exporter.Spans[0]

	attrs := attributesAsMap(sent.Resource().Attributes())

	assert.Equal(t, "test-span", sent.Name())
	assert.Equal(t, "some-value", attrs["service.namespace"].AsString())
}

func attributesAsMap(attrs []attribute.KeyValue) map[string]attribute.Value {
	m := make(map[string]attribute.Value, len(attrs))

	for _, attr := range attrs {
		m[string(attr.Key)] = attr.Value
	}

	return m
}
