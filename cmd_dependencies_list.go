package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func initDependenciesListCommand() *cobra.Command {
	var services *bool
	var displayJSON *bool
	var workingDirectory *string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all software that needs to be rebuilt or tested.",
		Long:  "Check all modified files from last commit and output the packages that needs to be rebuilt or tested.",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := modifiedPackages(cmd.Context(), *workingDirectory, *services)
			if err != nil {
				slog.Error("sp-build failed", slog.String("error", err.Error()))
				os.Exit(1)
			}
			displayList(list, *displayJSON)
		},
	}

	workingDirectory = cmd.Flags().StringP("working-directory", "w", ".", "working directory (directory that contains go.mod)")
	services = cmd.Flags().BoolP("services", "s", false, "only display services list")
	displayJSON = cmd.Flags().BoolP("json", "j", false, "display list as JSON")
	return cmd
}
