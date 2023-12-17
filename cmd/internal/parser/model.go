package parser

import "golang.org/x/tools/benchmark/parse"

type BenchmarkGroup struct {
	Name        string // Name of the dir where the benchmark is located
	Headline    string
	Description string
	Benchmarks  []Benchmark
	Code        string
	Constants   string
}

type Benchmark struct {
	Name        string // Name of the benchmark
	Description string // Description of the benchmark
	Code        string
	Variations  []Variation
}

type Variation struct {
	parse.Benchmark
	Name      string  // Name of the variation
	CPUCount  int     // Number of CPU cores used
	OpsPerSec float64 // Performance of the benchmark compared to the fastest benchmark
}

// --- BenchmarkMeta Model ---

type BenchmarkMeta struct {
	Name         string   `json:"name"`
	Headline     string   `json:"headline"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	Contributors []string `json:"contributors"`
	Meta         []struct {
		Implementation string `json:"implementation"`
		Description    string `json:"description"`
	} `json:"meta"`
}
