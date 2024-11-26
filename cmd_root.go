package main

import (
	"github.com/spf13/cobra"
)

func initRootCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sp-build",
		Short: "sp-build is a golang monorepo tool",
		Long:  `A Fast and Flexible golang monorepo tool`,
	}

	cmd.AddCommand(initDependenciesCommand())
	cmd.AddCommand(initServicesCommand())
	cmd.AddCommand(initVersionCommand(version))
	return cmd
}
