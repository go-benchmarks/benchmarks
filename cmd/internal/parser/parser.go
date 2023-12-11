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
		var set parse.Set

		benchmarkGroup.Name = filepath.Base(path)

		f, err := os.Open(path + string(os.PathSeparator) + "output.bench")
		if err != nil {
			return fmt.Errorf("failed to open benchmarkGroup file: %w", err)
		}

		set, err = parse.ParseSet(f)
		if err != nil {
			return fmt.Errorf("failed to parse benchmarkGroup file: %w", err)
		}

		var results []Benchmark
		for s, i := range set {
			logger.Debug("adding benchmarkGroup", "name", s)

			for _, b := range i {
				br := Benchmark{
					Benchmark: *b,
					Variation: s,
				}

				brNameParts := strings.Split(br.Name, "_")
				br.Name = brNameParts[0]
				if len(brNameParts) > 1 {
					br.Variation = strings.Join(brNameParts[1:], " ")
					br.Variation = strings.ReplaceAll(br.Variation, "_", " ")
					br.Variation = strings.ReplaceAll(br.Variation, "-", " ")

					// Variation parts.
					brVariationParts := strings.Split(br.Variation, " ")
					br.Variation = strings.Join(brVariationParts[:len(brVariationParts)-1], " ")

					// The last part is the CPU count, if it exists.
					br.CPUCount, err = strconv.Atoi(brVariationParts[len(brVariationParts)-1])
					if err != nil {
						br.CPUCount = 1
					}

				}

				// Split name. "BenchmarkName" -> "BenchmarkGroup Name". Split happens at every uppercase letter.
				br.Name = strings.Join(utils.SplitCamelCase(br.Name)[2:], " ")

				logger.Debug("adding benchmarkGroup", "name", br.Name, "variation", br.Variation, "cpuCount", br.CPUCount)

				results = append(results, br)
			}
		}

		// sort results by name and variation
		sort.Slice(results, func(i, j int) bool {
			if results[i].Name == results[j].Name {
				return results[i].Variation < results[j].Variation
			}

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
