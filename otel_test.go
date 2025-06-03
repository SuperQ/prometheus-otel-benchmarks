package main

import (
	"context"
	"log"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var (
	otelCounter      metric.Int64Counter
	ctx              = context.Background()
	cachedAttributes = attribute.NewSet(
		attribute.String("label1", "value1"),
		attribute.String("label2", "value2"),
	)
)

func init() {
	otlpMetricExporter, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL("http://localhost:9090/api/v1/otlp/v1/metrics"))
	if err != nil {
		log.Fatalf("Failed to create OTLP metric exporter: %v", err)
	}
	meterProvider := sdkmetric.NewMeterProvider(
		// Setting up an OTLP exporter is actually relevant for the benchmark, since increments flow
		// through to the exporter and drastically affect the performance.
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(otlpMetricExporter, sdkmetric.WithInterval(time.Hour))),
	)
	otel.SetMeterProvider(meterProvider)
	meter := otel.Meter("benchmark")

	otelCounter, err = meter.Int64Counter("test_counter",
		metric.WithDescription("A test counter"))
	if err != nil {
		panic(err)
	}
}

func BenchmarkOtelCounter(b *testing.B) {
	for b.Loop() {
		otelCounter.Add(ctx, 1)
	}
}

func BenchmarkOtelCounterWithAttributes(b *testing.B) {
	for b.Loop() {
		otelCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("label1", "value1"),
			attribute.String("label2", "value2")))
	}
}

func BenchmarkOtelCounterWithCachedAttributes(b *testing.B) {
	for b.Loop() {
		otelCounter.Add(ctx, 1, metric.WithAttributeSet(cachedAttributes))
	}
}

func BenchmarkOtelCounterParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			otelCounter.Add(ctx, 1)
		}
	})
}

func BenchmarkOtelCounterWithAttributesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			otelCounter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("label1", "value1"),
				attribute.String("label2", "value2"),
			))
		}
	})
}

func BenchmarkOtelCounterWithCachedAttributesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			otelCounter.Add(ctx, 1, metric.WithAttributeSet(cachedAttributes))
		}
	})
}
