package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func initServicesListCommand() *cobra.Command {
	var displayJSON *bool
	var goModDir *string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "command to list all services in the monorepo",
		Run: func(cmd *cobra.Command, args []string) {
			services, err := servicesList(cmd.Context(), *goModDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			displayMap(services, *displayJSON)
		},
	}

	goModDir = cmd.Flags().StringP("module", "m", ".", "path to the directory containing the go.mod file from repository root")
	displayJSON = cmd.Flags().BoolP("json", "j", false, "display list as JSON")
	return cmd
}
