package website

import (
	"bytes"
	"fmt"
	"github.com/go-benchmarks/benchmarks/cmd/internal/utils"
	"github.com/go-benchmarks/benchmarks/cmd/internal/website/templates"
	"golang.org/x/tools/benchmark/parse"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Website struct {
	Title      string
	Benchmarks []Benchmark
}

type Benchmark struct {
	Name               string // Name of the dir where the benchmark is located
	Results            []BenchmarkResult
	ResultsByVariation map[string][]BenchmarkResult
	ResultsByName      map[string][]BenchmarkResult
}

type BenchmarkResult struct {
	parse.Benchmark
	Variation string // Variation of the benchmark
	CPUCount  int    // Number of CPU cores used
}

func GenerateWebsite(logger *slog.Logger, benchmarksDir, outputDir string) error {
	website := Website{
		Title: "Go Benchmarks",
	}

	// Delete the output directory
	logger.Debug("deleting output directory", "outputDir", outputDir)
	err := os.RemoveAll(outputDir)
	if err != nil {
		return fmt.Errorf("failed to delete output directory: %w", err)
	}

	// Find all static pages
	var pages []string
	err = fs.WalkDir(templates.TemplateFS, "pages", func(path string, d fs.DirEntry, err error) error {
		logger.Debug("walking through TemplateFS", "currentPath", path)
		if err != nil {
			return fmt.Errorf("failed to walk through pages: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".gohtml" {
			return nil
		}

		pages = append(pages, path)

		return nil
	})
	if err != nil {
		return err
	}

	// Parse benchmarks
	err = utils.WalkOverBenchmarks(benchmarksDir, func(path string) error {
		logger.Debug("walking through benchmarks", "currentPath", path)

		var benchmark Benchmark
		var set parse.Set

		benchmark.Name = filepath.Base(path)

		f, err := os.Open(path + string(os.PathSeparator) + "output.bench")
		if err != nil {
			return fmt.Errorf("failed to open benchmark file: %w", err)
		}

		set, err = parse.ParseSet(f)
		if err != nil {
			return fmt.Errorf("failed to parse benchmark file: %w", err)
		}

		var results []BenchmarkResult
		for s, i := range set {
			logger.Debug("adding benchmark", "name", s)

			for _, b := range i {
				br := BenchmarkResult{
					Benchmark: *b,
					Variation: s,
				}

				brNameParts := strings.Split(br.Name, "_")
				br.Name = brNameParts[0]
				if len(brNameParts) > 1 {
					br.Variation = strings.Join(brNameParts[1:], " ")
					br.Variation = strings.ReplaceAll(br.Variation, "_", " ")
					br.Variation = strings.ReplaceAll(br.Variation, "-", " ")

					// Variation parts. The last part is the number of CPU cores used.
					brVariationParts := strings.Split(br.Variation, " ")
					br.Variation = strings.Join(brVariationParts[:len(brVariationParts)-1], " ")
					br.CPUCount, err = strconv.Atoi(brVariationParts[len(brVariationParts)-1])
					if err != nil {
						return fmt.Errorf("failed to parse CPU count: %w", err)
					}

				}

				// Split name. "BenchmarkName" -> "Benchmark Name". Split happens at every uppercase letter.
				br.Name = strings.Join(utils.SplitCamelCase(br.Name)[2:], " ")

				logger.Debug("adding benchmark", "name", br.Name, "variation", br.Variation, "cpuCount", br.CPUCount)

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

		// Results by variation
		benchmark.ResultsByVariation = make(map[string][]BenchmarkResult)
		for _, result := range results {
			benchmark.ResultsByVariation[result.Variation] = append(benchmark.ResultsByVariation[result.Variation], result)
		}

		// Results by name
		benchmark.ResultsByName = make(map[string][]BenchmarkResult)
		for _, result := range results {
			benchmark.ResultsByName[result.Name] = append(benchmark.ResultsByName[result.Name], result)
		}

		benchmark.Results = results

		website.Benchmarks = append(website.Benchmarks, benchmark)

		return nil
	})

	// Iterate over pages and generate the static pages with the base layout
	for _, page := range pages {
		logger.Debug("generating static page", "page", page)

		// Parse the page
		file, err := templates.TemplateFS.ReadFile(page)
		if err != nil {
			return err
		}

		tmpl := template.New("")
		_, err = tmpl.New(page).Parse(string(file))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		// Add the base layout
		_, err = tmpl.ParseFS(templates.TemplateFS, "layout/base.gohtml")
		if err != nil {
			return fmt.Errorf("failed to parse base layout: %w", err)
		}

		// Write the output file
		outputFilePath := filepath.Join(outputDir, page)
		outputFilePath = strings.ReplaceAll(outputFilePath, "pages"+string(os.PathSeparator), "")
		outputFilePath = strings.ReplaceAll(outputFilePath, ".gohtml", ".html")

		var output bytes.Buffer

		err = tmpl.ExecuteTemplate(&output, "base.gohtml", website)
		if err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		err = os.MkdirAll(filepath.Dir(outputFilePath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		logger.Debug("writing output file", "outputFilePath", outputFilePath)
		err = os.WriteFile(outputFilePath, output.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	// Render all benchmark pages
	for _, benchmark := range website.Benchmarks {
		tmpl := template.New("")

		// Add the base layout
		_, err = tmpl.ParseFS(templates.TemplateFS, "layout/base.gohtml")

		// Add "benchmark.gohtml"
		_, err = tmpl.ParseFS(templates.TemplateFS, "benchmark.gohtml")

		// Write the output file
		outputFilePath := filepath.Join(outputDir, "benchmark", benchmark.Name+string(os.PathSeparator)+"index.html")

		var output bytes.Buffer

		err = tmpl.ExecuteTemplate(&output, "base.gohtml", benchmark)
		if err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		err = os.MkdirAll(filepath.Dir(outputFilePath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		logger.Debug("writing output file", "outputFilePath", outputFilePath)
		err = os.WriteFile(outputFilePath, output.Bytes(), 0644)
	}

	return nil
}
