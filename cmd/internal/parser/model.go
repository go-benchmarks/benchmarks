package parser

import "golang.org/x/tools/benchmark/parse"

type BenchmarkGroup struct {
	Name       string      `json:"name,omitempty"` // Name of the dir where the benchmark is located
	Benchmarks []Benchmark `json:"benchmarks,omitempty"`
}

type Benchmark struct {
	parse.Benchmark `json:",inline"`
	Variation       string  `json:"variation,omitempty"`   // Variation of the benchmark
	CPUCount        int     `json:"CPUCount,omitempty"`    // Number of CPU cores used
	Performance     float64 `json:"performance,omitempty"` // Performance of the benchmark compared to the fastest benchmark
}
