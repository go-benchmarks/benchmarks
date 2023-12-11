package parser

import "golang.org/x/tools/benchmark/parse"

type BenchmarkGroup struct {
	Name        string // Name of the dir where the benchmark is located
	Description string
	Benchmarks  []Benchmark
}

type Benchmark struct {
	Name        string // Name of the benchmark
	Description string // Description of the benchmark
	Variations  []Variation
}

type Variation struct {
	parse.Benchmark
	Name      string  // Name of the variation
	CPUCount  int     // Number of CPU cores used
	OpsPerSec float64 // Performance of the benchmark compared to the fastest benchmark
}
