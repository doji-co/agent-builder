package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agent-builder",
	Short: "ADK Multi-Agent Builder CLI",
	Long:  "A CLI tool to help build ADK (Agent Development Kit) multi-agent systems.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
