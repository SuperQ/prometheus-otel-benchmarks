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

func init() {
	otlpMetricExporter, err := otlpmetrichttp.New(context.Background(), otlpmetrichttp.WithEndpointURL("http://localhost:9090/api/v1/otlp/v1/metrics"))
	if err != nil {
		log.Fatalf("Failed to create OTLP metric exporter: %v", err)
	}
	meterProvider := sdkmetric.NewMeterProvider(
		// Setting up an OTLP exporter is actually relevant for the benchmark, since increments flow
		// through to the exporter and drastically affect the performance.
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(otlpMetricExporter, sdkmetric.WithInterval(time.Hour))),
	)
	otel.SetMeterProvider(meterProvider)
}

func getOTelCounter() metric.Int64Counter {
	meter := otel.Meter("benchmark")
	otelCounter, err := meter.Int64Counter("test_counter",
		metric.WithDescription("A test counter"))
	if err != nil {
		panic(err)
	}
	return otelCounter
}

func getOTelHistogram() metric.Float64Histogram {
	meter := otel.Meter("benchmark")
	otelHistogram, err := meter.Float64Histogram("test_histogram",
		metric.WithDescription("A test histogram"))
	if err != nil {
		panic(err)
	}
	return otelHistogram
}

func BenchmarkOtelCounterParallel(b *testing.B) {
	ctx := b.Context()
	otelCounter := getOTelCounter()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			otelCounter.Add(ctx, 1)
		}
	})
}

func BenchmarkOtelCounterWithAttributesParallel(b *testing.B) {
	ctx := b.Context()
	otelCounter := getOTelCounter()

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
	ctx := b.Context()
	otelCounter := getOTelCounter()
	cachedAttributes := metric.WithAttributeSet(attribute.NewSet(
		attribute.String("label1", "value1"),
		attribute.String("label2", "value2"),
	))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			otelCounter.Add(ctx, 1, cachedAttributes)
		}
	})
}

func BenchmarkOtelHistogramParallel(b *testing.B) {
	ctx := b.Context()
	otelHistogram := getOTelHistogram()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			otelHistogram.Record(ctx, simulateObserve(b.N))
		}
	})
}

func BenchmarkOtelHistogramWithAttributesParallel(b *testing.B) {
	ctx := b.Context()
	otelHistogram := getOTelHistogram()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			otelHistogram.Record(ctx, simulateObserve(b.N), metric.WithAttributes(
				attribute.String("label1", "value1"),
				attribute.String("label2", "value2"),
			))
		}
	})
}

func BenchmarkOtelHistogramWithCachedAttributesParallel(b *testing.B) {
	ctx := b.Context()
	otelHistogram := getOTelHistogram()
	cachedAttributes := metric.WithAttributeSet(attribute.NewSet(
		attribute.String("label1", "value1"),
		attribute.String("label2", "value2"),
	))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			otelHistogram.Record(ctx, simulateObserve(b.N), cachedAttributes)
		}
	})
}
