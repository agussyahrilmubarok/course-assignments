package tracing

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	Tracer   trace.Tracer
	muTracer sync.RWMutex
)

func NewTracer(serviceName string) trace.Tracer {
	muTracer.Lock()
	defer muTracer.Unlock()

	Tracer = otel.Tracer(serviceName)
	return Tracer
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	muTracer.RLock()
	t := Tracer
	muTracer.RUnlock()

	if t == nil {
		t = otel.Tracer("default-tracer")
	}

	ctx, span := t.Start(ctx, name)
	return ctx, span
}
