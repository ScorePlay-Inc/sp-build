package main

import (
	"log/slog"
	"os"
)

var version = "unknown"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if err := initRootCommand(version).Execute(); err != nil {
		slog.Error("sp-build failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
