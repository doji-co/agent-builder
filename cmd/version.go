package cmd

import "github.com/spf13/cobra"

var buildVersion = "dev"

func SetVersion(version string) {
	if version == "" {
		buildVersion = "dev"
		rootCmd.Version = buildVersion
		return
	}
	buildVersion = version
	rootCmd.Version = buildVersion
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the current version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(buildVersion)
		},
	}
}
