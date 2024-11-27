package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

func commandContext(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd error: %w", err)
	}

	cmd := exec.CommandContext(ctx, name, arg...)
	slog.InfoContext(ctx, "running command",
		slog.String("name", cmd.String()),
		slog.String("path", currentPath),
	)
	return cmd, nil
}
