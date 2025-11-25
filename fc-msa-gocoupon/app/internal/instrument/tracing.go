package instrument

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type ShutdownFn func(ctx context.Context) error

func telemetryResource(ctx context.Context, serviceName string) *resource.Resource {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		log.Printf("failed to create telemetry resource %v", err)
	}
	return res
}

func InitTraceProvider(ctx context.Context, serviceName string, exporter sdktrace.SpanExporter) ShutdownFn {
	res := telemetryResource(ctx, serviceName)

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1.0)), // 100% sampling
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	return tp.Shutdown
}

func NewOTLPExporter(ctx context.Context, endpoint string) *otlptrace.Exporter {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(), // non TLS
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Printf("failed to create otlp trace exporter, %v", err)
	}
	return exporter
}

func NewZipkinExporter(endpoint string) sdktrace.SpanExporter {
	exporter, err := zipkin.New(endpoint)
	if err != nil {
		log.Fatalf("failed to create zipkin exporter: %v", err)
	}
	return exporter
}

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
