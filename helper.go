package main

import (
	"math"
)

func simulateObserve(i int) float64 {
	return 30 + math.Floor(120*math.Sin(float64(i)*0.1))/10
}
