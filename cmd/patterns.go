package cmd

import (
	"fmt"

	"github.com/doji-co/agent-builder/internal/model"
	"github.com/spf13/cobra"
)

func newPatternsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "patterns",
		Short: "List available orchestration patterns",
		Long:  "Display all available orchestration patterns with descriptions.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Available Orchestration Patterns:")
			for _, pattern := range []model.OrchestrationPattern{
				model.PatternSequential,
				model.PatternParallel,
				model.PatternLLMCoordinated,
				model.PatternLoop,
			} {
				cmd.Println(fmt.Sprintf("%s: %s", pattern.String(), pattern.Description()))
			}
		},
	}
}
