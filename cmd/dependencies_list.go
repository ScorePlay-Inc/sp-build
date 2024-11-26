package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ScorePlay-Inc/media-management/tools/revdeps/internal/usecases/dependencies"
	"github.com/ScorePlay-Inc/media-management/tools/revdeps/internal/utils"
)

func initDependenciesListCommand() *cobra.Command {
	var services *bool
	var displayJSON *bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all software that needs to be rebuilt or tested.",
		Long:  "Check all modified files from last commit and output the packages that needs to be rebuilt or tested.",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := dependencies.ModifiedPackages(cmd.Context(), *services)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			utils.DisplayList(list, *displayJSON)
		},
	}

	services = cmd.Flags().BoolP("services", "s", false, "only display services list")
	displayJSON = cmd.Flags().BoolP("json", "j", false, "display list as JSON")
	return cmd
}
