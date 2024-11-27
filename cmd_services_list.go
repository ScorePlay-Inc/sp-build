package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func initServicesListCommand() *cobra.Command {
	var displayJSON *bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "command to list all services in the monorepo",
		Run: func(cmd *cobra.Command, args []string) {
			services, err := servicesList(cmd.Context())
			if err != nil {
				slog.Error("sp-build failed", slog.String("error", err.Error()))
				os.Exit(1)
			}
			displayMap(services, *displayJSON)
		},
	}

	displayJSON = cmd.Flags().BoolP("json", "j", false, "display list as JSON")
	return cmd
}
