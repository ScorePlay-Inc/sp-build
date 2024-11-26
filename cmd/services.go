package cmd

import (
	"github.com/spf13/cobra"
)

func initServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "command line to help with the monorepo dependencies",
	}

	cmd.AddCommand(initServicesListCommand())
	return cmd
}
