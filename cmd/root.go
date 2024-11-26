package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func initRootCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revdeps",
		Short: "revdeps is a golang monorepo tool",
		Long:  `A Fast and Flexible golang monorepo tool`,
	}

	cmd.AddCommand(initDependenciesCommand())
	cmd.AddCommand(initServicesCommand())
	cmd.AddCommand(initVersionCommand(version))
	return cmd
}

func Execute(version string) {
	if err := initRootCommand(version).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
