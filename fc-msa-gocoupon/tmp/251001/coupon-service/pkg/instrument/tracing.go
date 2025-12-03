package instrument

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
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
		log.Fatal().Err(err).Msg("failed to create telemetry resource")
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
		otlptracegrpc.WithInsecure(),         // non TLS
		otlptracegrpc.WithEndpoint(endpoint), // e.g: "tempo:4317"
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create OTLP trace exporter")
	}
	return exporter
}
