package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
)

func commandContext(ctx context.Context, workingDirectory string, name string, arg ...string) (*exec.Cmd, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd error: %w", err)
	}

	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Dir = workingDirectory

	slog.InfoContext(ctx, "running command",
		slog.String("name", cmd.String()),
		slog.String("path", path.Join(currentPath, workingDirectory)),
	)
	return cmd, nil
}
