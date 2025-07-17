package main

// Benchmark how fast we can increment a Prometheus counter metric without any labels
// and another one with two labels.
import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func BenchmarkPrometheusCounterParallel(b *testing.B) {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "A test counter",
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Inc()
		}
	})
}

func BenchmarkPrometheusCounterWithLabelsParallel(b *testing.B) {
	counterWithLabels := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "A test counter",
	}, []string{"label1", "label2"})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counterWithLabels.WithLabelValues("value1", "value2").Inc()
		}
	})
}

func BenchmarkPrometheusCounterWithCachedLabelsParallel(b *testing.B) {
	counterWithCachedLabels := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "A test counter",
	}, []string{"label1", "label2"}).WithLabelValues("value1", "value2")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counterWithCachedLabels.Inc()
		}
	})
}
