package otel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func InitTracer(collectorEndpoint string, serviceName string) {
	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(collectorEndpoint),
		otlptracehttp.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("Error creating OTLP exporter: %v", err)
	}

	resourceName, err := resource.New(
		nil,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)

	if err != nil {
		log.Fatalf("Error creating resource: %v", err)
	}

	batchSpan := trace.NewBatchSpanProcessor(exporter)

	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(resourceName),
		trace.WithSpanProcessor(batchSpan),
	)

	otel.SetTracerProvider(traceProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})
}
