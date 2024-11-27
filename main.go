package main

import (
	"fmt"
	"log/slog"
	"os"
)

var version = "unknown"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if err := initRootCommand(version).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
