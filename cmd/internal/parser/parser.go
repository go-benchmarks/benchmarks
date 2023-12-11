package parser

import (
	"encoding/json"
	"fmt"
	"github.com/go-benchmarks/benchmarks/cmd/internal/utils"
	"golang.org/x/tools/benchmark/parse"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func GenerateJson(logger *slog.Logger, benchmarksDir string, pretty bool) ([]byte, error) {
	benchmarkGroups, err := ProcessBenchmarkGroups(logger, benchmarksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to process benchmark groups: %w", err)
	}

	if pretty {
		return json.MarshalIndent(benchmarkGroups, "", "  ")
	}

	return json.Marshal(benchmarkGroups)
}

func ProcessBenchmarkGroups(logger *slog.Logger, benchmarksDir string) (groups []BenchmarkGroup, err error) {
	// Parse benchmarks
	err = utils.WalkOverBenchmarks(benchmarksDir, func(path string) error {
		logger.Debug("walking through benchmarks", "currentPath", path)

		var benchmarkGroup BenchmarkGroup

		benchmarkGroup.Name = filepath.Base(path)

		f, err := os.Open(path + string(os.PathSeparator) + "output.bench")
		if err != nil {
			return fmt.Errorf("failed to open benchmarkGroup file: %w", err)
		}

		set, err := parse.ParseSet(f)
		if err != nil {
			return fmt.Errorf("failed to parse benchmarkGroup file: %w", err)
		}

		var variations []Variation
		for s, i := range set {
			for _, b := range i {
				logger.Debug("adding variation", "name", s)
				variation := Variation{
					Benchmark: *b,
				}

				brNameParts := strings.Split(variation.Benchmark.Name, "_") // "BenchmarkName_VariationName" -> ["BenchmarkName", "VariationName"]
				logger.Debug("benchmark name parts", "parts", brNameParts)
				variation.Benchmark.Name = brNameParts[0] // Benchmark name is the first part.

				// If there are more parts, then the variation name is the second part.
				if len(brNameParts) > 1 {
					variation.Name = brNameParts[1]
					variation.Name = strings.ReplaceAll(variation.Name, "_", " ")
					variation.Name = strings.ReplaceAll(variation.Name, "-", " ")

					// Variation parts.
					brVariationParts := strings.Split(variation.Name, " ")
					logger.Debug("benchmark variation name parts", "parts", brVariationParts)
					variation.Name = strings.Join(brVariationParts[:len(brVariationParts)-1], " ")

					// The last part is the CPU count, if it exists.
					variation.CPUCount, err = strconv.Atoi(brVariationParts[len(brVariationParts)-1])
					if err != nil {
						variation.CPUCount = 1
						variation.Name = strings.Join(brVariationParts, " ")
					}
				}

				// Split name. "BenchmarkName" -> "BenchmarkGroup Name". Split happens at every uppercase letter.
				variation.Benchmark.Name = strings.Join(utils.SplitCamelCase(variation.Benchmark.Name)[2:], " ")
				logger.Debug("adding benchmark variation", "benchmark name", variation.Benchmark.Name, "variation name", variation.Name, "cpuCount", variation.CPUCount, "orig name", s)

				// Calculate ops per second by dividing ns/op by 1e9.
				variation.OpsPerSec = 1e9 / variation.NsPerOp

				variations = append(variations, variation)
			}
		}

		logger.Info("added benchmark variations", "count", len(variations))

		benchmarks := make(map[string][]Variation)
		for _, v := range variations {
			benchmarks[v.Benchmark.Name] = append(benchmarks[v.Benchmark.Name], v)
		}

		var results []Benchmark
		for name, variations := range benchmarks {
			var benchmark Benchmark
			benchmark.Name = name
			benchmark.Variations = variations

			results = append(results, benchmark)
		}

		// sort results by name
		sort.Slice(results, func(i, j int) bool {
			return results[i].Name < results[j].Name
		})

		benchmarkGroup.Benchmarks = results

		groups = append(groups, benchmarkGroup)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return groups, nil
}
