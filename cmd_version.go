package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func initVersionCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "display sp-build version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	return cmd
}
