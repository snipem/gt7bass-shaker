package main

import "testing"

func Benchmark_shift(b *testing.B) {

	for i := 0; i < 1000; i++ {
		shift()
	}
}
