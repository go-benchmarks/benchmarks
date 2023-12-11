package parser

import "golang.org/x/tools/benchmark/parse"

type BenchmarkGroup struct {
	Name       string // Name of the dir where the benchmark is located
	Benchmarks []Benchmark
}

type Benchmark struct {
	parse.Benchmark
	Variation   string  // Variation of the benchmark
	CPUCount    int     // Number of CPU cores used
	Performance float64 // Performance of the benchmark compared to the fastest benchmark
}
