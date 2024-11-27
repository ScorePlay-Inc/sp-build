package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"path"
	"strings"
)

func golangModuleName(ctx context.Context, workingDirectory string) (string, error) {
	cmd, err := commandContext(ctx, workingDirectory, "go", "mod", "edit", "--json")
	if err != nil {
		return "", fmt.Errorf("commandContext error: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cmd.CombinedOutput error (%w): %s", err, string(output))
	}

	var module struct {
		Module struct {
			Path string `json:"Path"`
		} `json:"Module"`
	}

	if err := json.Unmarshal(output, &module); err != nil {
		return "", fmt.Errorf("json.Unmarshal error: %w", err)
	}

	name := module.Module.Path
	slog.InfoContext(ctx, "fetching golang module name", slog.String("name", name))
	return name, nil
}

func modifiedFilesSinceLastCommit(ctx context.Context, workingDirectory string) ([]string, error) {
	listing, err := exec.CommandContext(ctx, "ls", "-a").Output()
	if err != nil {
		return nil, fmt.Errorf("ls -la error: %w", err)
	}

	content := strings.Split(strings.TrimSpace(string(listing)), "\n")
	slog.InfoContext(ctx, "listing root directory content",
		slog.String("listing", strings.Join(content, " ")),
	)

	gitStatus, err := exec.CommandContext(ctx, "git", "status", "--porcelain").Output()
	if err != nil {
		return nil, fmt.Errorf("git status error: %w", err)
	}

	status := strings.Split(strings.TrimSpace(string(gitStatus)), "\n")
	slog.InfoContext(ctx, "fetching status of the repository",
		slog.String("status", strings.Join(status, " ")),
	)

	cmd, err := commandContext(ctx, workingDirectory, "git", "diff", "--relative", "--name-only", "@^")
	if err != nil {
		return nil, fmt.Errorf("commandContext error: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cmd.CombinedOutput error (%w): %s", err, string(output))
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	slog.InfoContext(ctx, "fetching modified files",
		slog.String("files", strings.Join(files, ", ")),
		slog.Int("count", len(files)),
	)

	return files, nil
}

func modifiedPackages(ctx context.Context, workingDirectory string, onlyServices bool) ([]string, error) {
	moduleName, err := golangModuleName(ctx, workingDirectory)
	if err != nil {
		return nil, fmt.Errorf("golangModuleName error %w", err)
	}

	modifiedFiles, err := modifiedFilesSinceLastCommit(ctx, workingDirectory)
	if err != nil {
		return nil, fmt.Errorf("modifiedFilesSinceLastCommit error %w", err)
	}

	revDeps, err := getReverseDependencies(ctx, workingDirectory, onlyServices)
	if err != nil {
		return nil, fmt.Errorf("getReverseDependencies error %w", err)
	}

	modified := make(map[string]bool)
	for _, file := range modifiedFiles {
		file = path.Join(moduleName, path.Dir(file))

		for _, pkg := range revDeps[file] {
			modified[pkg] = true
		}
	}

	packageList := make([]string, 0)
	for pkg := range modified {
		pkgPath := strings.TrimPrefix(pkg, moduleName)
		if pkgPath == "" {
			continue
		}

		packageList = append(packageList, "."+pkgPath)
	}
	return packageList, nil
}

func servicesList(ctx context.Context, workingDirectory string) (map[string]string, error) {
	goModuleName, err := golangModuleName(ctx, workingDirectory)
	if err != nil {
		return nil, fmt.Errorf("golangModuleName error %w", err)
	}

	services, err := getServicesList(ctx, workingDirectory, goModuleName)
	if err != nil {
		return nil, fmt.Errorf("getServicesList error: %w", err)
	}
	return services, nil
}
