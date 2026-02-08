package tracing

import (
	"context"
	"fmt"
	"time"

	"kasir-api/pkg/logger"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// TraceRequest traces a request with start and end times
func TraceRequest(ctx context.Context, operation string, req interface{}) (context.Context, func(resp interface{}, err error)) {
	start := time.Now()
	requestID := fmt.Sprintf("%d", time.Now().UnixNano()) // Simple ID generation

	ctx = context.WithValue(ctx, requestIDKey, requestID)

	logger.InfoCtx(ctx, "Request started",
		"operation", operation,
		"timestamp", start.Format(time.RFC3339),
		"request", req,
	)

	return ctx, func(resp interface{}, err error) {
		duration := time.Since(start)

		if err != nil {
			logger.ErrorCtx(ctx, "Request completed with error",
				"operation", operation,
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
				"response", resp,
			)
		} else {
			logger.InfoCtx(ctx, "Request completed",
				"operation", operation,
				"duration_ms", duration.Milliseconds(),
				"response_size", estimateResponseSize(resp),
				"response", resp,
			)
		}
	}
}

// estimateResponseSize provides a rough estimate of response size
func estimateResponseSize(resp interface{}) int {
	if resp == nil {
		return 0
	}

	// This is a simple estimation - in a real system you might want to use
	// reflection or serialization to get a more accurate size
	switch v := resp.(type) {
	case string:
		return len(v)
	case []byte:
		return len(v)
	default:
		return 100 // default estimate
	}
}

// Span represents a tracing span
type Span struct {
	Operation string
	Start     time.Time
	Context   context.Context
}

// NewSpan creates a new tracing span
func NewSpan(ctx context.Context, operation string) *Span {
	return &Span{
		Operation: operation,
		Start:     time.Now(),
		Context:   ctx,
	}
}

// End ends the span and logs the result
func (s *Span) End(result interface{}, err error) {
	duration := time.Since(s.Start)

	if err != nil {
		logger.ErrorCtx(s.Context, "Span completed with error",
			"operation", s.Operation,
			"duration_ms", duration.Milliseconds(),
			"error", err.Error(),
		)
	} else {
		logger.InfoCtx(s.Context, "Span completed",
			"operation", s.Operation,
			"duration_ms", duration.Milliseconds(),
		)
	}
}

// WithAttribute adds attributes to the span
func (s *Span) WithAttribute(key string, value interface{}) *Span {
	// In a real tracing system, you would add attributes to the span
	// For now, we'll just return the same span
	return s
}
