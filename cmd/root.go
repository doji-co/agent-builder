package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = newRootCommand()

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent-builder",
		Short: "ADK Multi-Agent Builder CLI",
		Long:  "A CLI tool to help build ADK multi-agent systems.",
	}
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(
		newCreateCommand(createDependencies{}),
		newPatternsCommand(),
		newVersionCommand(),
	)
	return cmd
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
