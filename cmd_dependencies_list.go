package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func initDependenciesListCommand() *cobra.Command {
	var services *bool
	var displayJSON *bool
	var goModDir *string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all software that needs to be rebuilt or tested.",
		Long:  "Check all modified files from last commit and output the packages that needs to be rebuilt or tested.",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := modifiedPackages(cmd.Context(), *goModDir, *services)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			displayList(list, *displayJSON)
		},
	}

	goModDir = cmd.Flags().StringP("module", "m", "./", "path to the directory containing the go.mod file")
	services = cmd.Flags().BoolP("services", "s", false, "only display services list")
	displayJSON = cmd.Flags().BoolP("json", "j", false, "display list as JSON")
	return cmd
}
