package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ScorePlay-Inc/media-management/tools/revdeps/internal/usecases/dependencies"
	"github.com/ScorePlay-Inc/media-management/tools/revdeps/internal/utils"
)

func initServicesListCommand() *cobra.Command {
	var displayJSON *bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "command to list all services in the monorepo",
		Run: func(cmd *cobra.Command, args []string) {
			services, err := dependencies.ServicesList(cmd.Context())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			utils.DisplayMap(services, *displayJSON)
		},
	}

	displayJSON = cmd.Flags().BoolP("json", "j", false, "display list as JSON")
	return cmd
}
