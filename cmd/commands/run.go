package commands

import (
	"fmt"
	"github.com/go-benchmarks/benchmarks/cmd/internal/utils"
	"github.com/go-benchmarks/benchmarks/cmd/logger"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run benchmarks and save results",
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		logger := logger.New(debug)

		logger.Info("running benchmarks")

		basePath := cmd.Flag("benchmarks").Value.String()
		logger.Debug("benchmarks directory", "basePath", basePath)
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			return fmt.Errorf("benchmarks directory does not exist: %s", basePath)
		}

		// Walk through the benchmarks directory
		err := utils.WalkOverBenchmarks(basePath, func(path string) error {
			return runBenchmark(logger, path)
		})

		return err
	},
}

func runBenchmark(logger *slog.Logger, path string) error {
	logger.Debug("running benchmark", "path", path)

	maxCPU := runtime.NumCPU()
	logger.Debug("max cpu", "maxCPU", maxCPU)

	var cpuTests []string
	for i := 1; i <= maxCPU; i *= 2 {
		cpuTests = append(cpuTests, fmt.Sprint(i))
	}

	logger.Debug("cpu tests", "cpuTests", cpuTests)

	cmd := exec.Command("go", "test", "-bench", ".", "-benchmem", "-cpu", strings.Join(cpuTests, ","))
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("failed to run benchmark", "path", path, "output", string(output))
		return fmt.Errorf("failed to run benchmark: %w", err)
	}

	logger.Info("writing benchmark output", "path", path+string(os.PathSeparator)+"output.bench")
	outputFilePath := filepath.Join(path, "output.bench")
	return os.WriteFile(outputFilePath, output, 0644)
}

func init() {
	runCmd.Flags().StringP("benchmarks", "b", "./benchmarks", "Filepath of the \"benchmarks\" directory")

	rootCmd.AddCommand(runCmd)
}
