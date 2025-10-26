package cmd

import (
	"fmt"

	"github.com/doji-co/agent-builder/internal/model"
	"github.com/spf13/cobra"
)

var patternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "List available orchestration patterns",
	Long:  "Display all available orchestration patterns with descriptions.",
	Run:   runPatterns,
}

func init() {
	rootCmd.AddCommand(patternsCmd)
}

func runPatterns(cmd *cobra.Command, args []string) {
	fmt.Println("Available Orchestration Patterns:")

	patterns := []model.OrchestrationPattern{
		model.PatternSequential,
		model.PatternParallel,
		model.PatternLLMCoordinated,
		model.PatternLoop,
	}

	for _, pattern := range patterns {
		fmt.Printf("• %s\n", pattern.String())
		fmt.Printf("  %s\n\n", pattern.Description())
	}
}
