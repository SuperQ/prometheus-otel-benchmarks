package main

// Benchmark how fast we can increment a Prometheus counter metric without any labels
// and another one with two labels.
import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "A test counter",
	})

	counterWithLabels = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "test_counter_with_labels",
		Help: "A test counter with labels",
	}, []string{"label1", "label2"})

	counterWithCachedLabels = counterWithLabels.WithLabelValues("value1", "value2")
)

func init() {
	prometheus.MustRegister(counter)
}
func BenchmarkPrometheusCounter(b *testing.B) {
	for b.Loop() {
		counter.Inc()
	}
}

func BenchmarkPrometheusCounterWithLabels(b *testing.B) {
	for b.Loop() {
		counterWithLabels.WithLabelValues("value1", "value2").Inc()
	}
}

func BenchmarkPrometheusCounterWithCachedLabels(b *testing.B) {
	for b.Loop() {
		counterWithCachedLabels.Inc()
	}
}

func BenchmarkPrometheusCounterParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Inc()
		}
	})
}

func BenchmarkPrometheusCounterWithLabelsParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counterWithLabels.WithLabelValues("value1", "value2").Inc()
		}
	})
}

func BenchmarkPrometheusCounterWithCachedLabelsParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counterWithCachedLabels.Inc()
		}
	})
}
