package commands

import (
	"fmt"
	"github.com/go-benchmarks/benchmarks/cmd/internal/parser"
	"github.com/go-benchmarks/benchmarks/cmd/logger"
	"github.com/spf13/cobra"
	"os"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate benchmarks",
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		logger := logger.New(debug)

		benchmarksDir := cmd.Flag("benchmarks").Value.String()
		filePath := cmd.Flag("file").Value.String()
		logger.Debug("flags", "benchmarkDir", benchmarksDir, "outputDir", filePath)

		logger.Debug("checking if benchmarks directory exists")
		if _, err := os.Stat(benchmarksDir); os.IsNotExist(err) {
			return fmt.Errorf("benchmarks directory does not exist: %s", benchmarksDir)
		}

		j, err := parser.GenerateJson(logger, benchmarksDir, false)
		if err != nil {
			return fmt.Errorf("failed to generate output: %w", err)
		}

		// Write json file
		logger.Info("writing json file", "path", filePath)
		err = os.WriteFile(filePath, j, 0644)

		return nil
	},
}

func init() {
	generateCmd.Flags().StringP("benchmarks", "b", "./benchmarks", "Filepath of the \"benchmarks\" directory")
	generateCmd.Flags().StringP("file", "f", "./web/src/lib/benchmarks.json", "Path of the output file")

	rootCmd.AddCommand(generateCmd)
}
