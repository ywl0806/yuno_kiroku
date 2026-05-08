package utils

import (
	"fmt"
	"strings"
)

// Float64SliceToVectorString converts []float64 to pgvector string format
// Example: [0.1, 0.2, 0.3] -> "[0.1,0.2,0.3]"
func Float64SliceToVectorString(embedding []float64) string {
	if len(embedding) == 0 {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteString("[")
	for i, v := range embedding {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("%g", v))
	}
	builder.WriteString("]")
	return builder.String()
}
