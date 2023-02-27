package tracing

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

func AttributesFromMap[V any](m map[string]V) []attribute.KeyValue {

	attrs := make([]attribute.KeyValue, 0, len(m))

	for k, v := range m {
		if strings.HasPrefix(k, "attr.") {
			attrs = append(attrs, asAttribute(strings.TrimPrefix(k, "attr."), v))
		}
	}

	return attrs
}

func asAttribute(key string, v any) attribute.KeyValue {

	switch val := v.(type) {

	case string:
		return attribute.String(key, val)

	case *string:
		return attribute.String(key, *val)

	case bool:
		return attribute.Bool(key, val)

	case *bool:
		return attribute.Bool(key, *val)

	case int:
		return attribute.Int(key, val)

	case *int:
		return attribute.Int(key, *val)

	case int64:
		return attribute.Int64(key, val)

	case *int64:
		return attribute.Int64(key, *val)

	case *time.Time:
		return attribute.Int64(key, val.UnixMilli())

	case time.Time:
		return attribute.Int64(key, val.UnixMilli())

	default:
		// handle the pointer case, so that `Sprint` prints the actual value
		if reflect.ValueOf(v).Kind() == reflect.Pointer {
			v = reflect.Indirect(reflect.ValueOf(v))
		}

		return attribute.String(key, fmt.Sprint(v))

	}
}
