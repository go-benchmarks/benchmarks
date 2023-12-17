package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/go-benchmarks/benchmarks/cmd/internal/utils"
	"github.com/goccy/go-yaml"
	"go/token"
	"golang.org/x/tools/benchmark/parse"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
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

		f, err := os.Open(path + string(os.PathSeparator) + "output.bench")
		if err != nil {
			return fmt.Errorf("failed to open benchmarkGroup file: %w", err)
		}

		set, err := parse.ParseSet(f)
		if err != nil {
			return fmt.Errorf("failed to parse benchmarkGroup file: %w", err)
		}

		// Init BenchmarkMeta
		var meta BenchmarkMeta
		meta.Name = filepath.Base(path)

		// Check if _meta.yaml exists
		metaFilePath := path + string(os.PathSeparator) + "_meta.yaml"
		if _, err := os.Stat(metaFilePath); err == nil {
			logger.Debug("meta file exists", "path", metaFilePath)

			// Open meta file
			metaFile, err := os.Open(metaFilePath)
			if err != nil {
				return fmt.Errorf("failed to open meta file: %w", err)
			}

			// Decode meta file
			err = yaml.NewDecoder(metaFile).Decode(&meta)
			if err != nil {
				return fmt.Errorf("failed to decode meta file: %w", err)
			}
		} else {
			logger.Warn("no meta file found", "path", metaFilePath)
		}

		// Get all *_test.go files
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".go") {
				logger.Debug("found test file", "path", path)

				// Read test file
				b, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read test file: %w", err)
				}

				cC, err := cleanCode(string(b))
				if err != nil {
					return fmt.Errorf("failed to clean test file: %w", err)
				}

				benchmarkGroup.Code += cC

				consts, err := getConsts(string(b))
				if err != nil {
					return fmt.Errorf("failed to get consts: %w", err)
				}

				benchmarkGroup.Constants += consts
			}

			return nil
		})

		benchmarkGroup.Code = strings.TrimSpace(benchmarkGroup.Code)

		benchmarkGroup.Name = meta.Name
		benchmarkGroup.Description = meta.Description
		benchmarkGroup.Headline = meta.Headline

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
				variation.Benchmark.Name = strings.Join(utils.SplitCamelCase(variation.Benchmark.Name)[1:], " ")
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

			for _, m := range meta.Meta {
				if m.Implementation == name {
					benchmark.Description = m.Description
				}
			}

			logger.Debug("getting code", "benchmark name", name)
			benchmark.Code, err = getCode(benchmarkGroup.Code, strings.ReplaceAll(name, " ", ""))
			if err != nil {
				return fmt.Errorf("failed to get benchmark code: %w", err)
			}

			logger.Debug("getting benchmark code", "benchmark name", name)
			benchmark.BenchmarkCode, err = getBenchmarkCode(benchmarkGroup.Code, strings.ReplaceAll(name, " ", ""))
			if err != nil {
				return fmt.Errorf("failed to get benchmark code: %w", err)
			}

			benchmark.Code = strings.TrimSpace(benchmark.Code)
			benchmark.BenchmarkCode = strings.TrimSpace(benchmark.BenchmarkCode)

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

func cleanCode(src string) (string, error) {
	// Remove import blocks and lines that start with "package"
	re := regexp.MustCompile(`(?m)^import \([\s\S]*?\)\n|^import .*\n|^package .*\n`)
	src = re.ReplaceAllString(src, "")

	// Replace multiple consecutive newline characters with a single newline character
	re = regexp.MustCompile(`\n{3,}`)
	src = re.ReplaceAllString(src, "\n")

	src = strings.TrimSpace(src)
	src += "\n\n"

	if src == "\n\n" {
		src = ""
	}

	return src, nil
}

func getConsts(src string) (string, error) {
	file, err := decorator.Parse(src)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	newFile := &dst.File{}
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *dst.GenDecl:
			if decl.Tok == token.CONST {
				newFile.Decls = append(newFile.Decls, decl)
			}
		}
	}

	newFile.Name = dst.NewIdent("dummy")
	decorator.Fprint(&buf, newFile)
	return cleanCode(buf.String())
}

func getBenchmarkCode(src, name string) (string, error) {
	src, _ = cleanCode(src)
	src = "package dummy\n\n" + src
	file, err := decorator.Parse(src)
	if err != nil {
		return "", err
	}

	if file == nil {
		return "", fmt.Errorf("parsed file is nil")
	}

	var buf bytes.Buffer
	newFile := &dst.File{}
	dst.Inspect(file, func(n dst.Node) bool {
		switch decl := n.(type) {
		case *dst.FuncDecl:
			if strings.HasPrefix(decl.Name.Name, "Benchmark"+name+"_") {
				newFile.Decls = append(newFile.Decls, decl)
			}
		}
		return true
	})

	newFile.Name = dst.NewIdent("dummy")
	decorator.Fprint(&buf, newFile)
	return cleanCode(buf.String())
}

func getCode(src, name string) (string, error) {
	src, _ = cleanCode(src)
	src = "package dummy\n\n" + src
	file, err := decorator.Parse(src)
	if err != nil {
		return "", err
	}

	if file == nil {
		return "", fmt.Errorf("parsed file is nil")
	}

	var buf bytes.Buffer
	newFile := &dst.File{}
	dst.Inspect(file, func(n dst.Node) bool {
		switch decl := n.(type) {
		case *dst.GenDecl:
			if decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					typeSpec, ok := spec.(*dst.TypeSpec)
					if ok && typeSpec.Name.Name == name {
						newFile.Decls = append(newFile.Decls, decl)
					}
				}
			}
		case *dst.FuncDecl:
			if decl.Recv != nil && len(decl.Recv.List) > 0 {
				var recvTypeName string

				switch expr := decl.Recv.List[0].Type.(type) {
				case *dst.StarExpr:
					if expr.X != nil {
						recvTypeName = expr.X.(*dst.Ident).Name
					}
				case *dst.Ident:
					recvTypeName = expr.Name
				}

				if recvTypeName == name {
					newFile.Decls = append(newFile.Decls, decl)
				}
			}
		}
		return true
	})

	newFile.Name = dst.NewIdent("dummy")
	decorator.Fprint(&buf, newFile)
	return cleanCode(buf.String())
}
