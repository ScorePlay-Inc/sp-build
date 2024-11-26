package cmd

import (
	"github.com/spf13/cobra"
)

func initDependenciesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dependencies",
		Short: "command line to help with the monorepo dependencies",
	}

	cmd.AddCommand(initDependenciesListCommand())
	return cmd
}
