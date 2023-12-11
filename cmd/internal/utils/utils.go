package utils

import (
	"os"
	"path/filepath"
)

func WalkOverBenchmarks(basePath string, f func(path string) error) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the base benchmarks directory
		if path == basePath {
			return nil
		}

		if info.IsDir() {
			return f(path)
		}

		return nil
	})
}

func SplitCamelCase(input string) []string {
	var output []string

	var current string
	for _, r := range input {
		if r >= 'A' && r <= 'Z' {
			output = append(output, current)
			current = ""
		}

		current += string(r)
	}

	if current != "" {
		output = append(output, current)
	}

	return output
}
