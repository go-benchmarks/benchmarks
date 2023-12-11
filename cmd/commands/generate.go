package commands

import (
	"fmt"
	"github.com/go-benchmarks/benchmarks/cmd/internal/website"
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
		outputDir := cmd.Flag("output").Value.String()
		logger.Debug("flags", "benchmarkDir", benchmarksDir, "outputDir", outputDir)

		logger.Debug("checking if benchmarks directory exists")
		if _, err := os.Stat(benchmarksDir); os.IsNotExist(err) {
			return fmt.Errorf("benchmarks directory does not exist: %s", benchmarksDir)
		}

		logger.Debug("checking if output directory exists")
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			// Create the output directory
			logger.Info("creating output directory", "path", outputDir)
			err := os.MkdirAll(outputDir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
		}

		err := website.GenerateWebsite(logger, benchmarksDir, outputDir)
		if err != nil {
			return fmt.Errorf("failed to generate output: %w", err)
		}

		return nil
	},
}

func init() {
	generateCmd.Flags().StringP("benchmarks", "b", "./benchmarks", "Filepath of the \"benchmarks\" directory")
	generateCmd.Flags().StringP("output", "o", "./docs", "Filepath of the output directory")

	rootCmd.AddCommand(generateCmd)
}
