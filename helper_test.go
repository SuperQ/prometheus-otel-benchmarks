package main

import (
	"testing"
)

func BenchmarkSimulateObserve(b *testing.B) {
	for b.Loop() {
		_ = simulateObserve(b.N)
	}
}
