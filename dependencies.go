package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func golangModuleName(ctx context.Context, goModDir string) (string, error) {
	cmd := exec.CommandContext(ctx, "go", "mod", "edit", "--json")
	cmd.Dir = goModDir

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("exec.Command error: %w", err)
	}

	var module struct {
		Module struct {
			Path string `json:"Path"`
		} `json:"Module"`
	}

	if err := json.Unmarshal(output, &module); err != nil {
		return "", fmt.Errorf("json.Unmarshal error: %w", err)
	}
	return module.Module.Path, nil
}

func modifiedFilesSinceLastCommit(ctx context.Context) ([]string, error) {
	modifiedFiles, err := exec.CommandContext(ctx, "git", "diff", "--name-only", "@^").Output()
	if err != nil {
		return nil, fmt.Errorf("exec.Command error: %w", err)
	}
	return strings.Split(strings.TrimSpace(string(modifiedFiles)), "\n"), nil
}

func modifiedPackages(ctx context.Context, goModDir string, onlyServices bool) ([]string, error) {
	moduleName, err := golangModuleName(ctx, goModDir)
	if err != nil {
		return nil, fmt.Errorf("golangModuleName error %w", err)
	}

	modifiedFiles, err := modifiedFilesSinceLastCommit(ctx)
	if err != nil {
		return nil, fmt.Errorf("modifiedFilesSinceLastCommit error %w", err)
	}

	revDeps, err := getReverseDependencies(ctx, goModDir, onlyServices)
	if err != nil {
		return nil, fmt.Errorf("getReverseDependencies error %w", err)
	}

	toBuild := make(map[string]bool)
	for _, file := range modifiedFiles {
		if strings.HasPrefix(file, "vendor/") {
			file = strings.TrimPrefix(path.Dir(file), "vendor/")
		} else {
			file = path.Join(moduleName, path.Dir(file))
		}

		for _, pkg := range revDeps[file] {
			toBuild[pkg] = true
		}
	}

	packageList := make([]string, 0)
	for service := range toBuild {
		serviceName := strings.TrimPrefix(service, moduleName)
		if serviceName == "" {
			continue
		}

		packageList = append(packageList, serviceName)
	}
	return packageList, nil
}

func servicesList(ctx context.Context, goModDir string) (map[string]string, error) {
	goModuleName, err := golangModuleName(ctx, goModDir)
	if err != nil {
		return nil, fmt.Errorf("golangModuleName error %w", err)
	}

	services, err := getServicesList(ctx, goModDir, goModuleName)
	if err != nil {
		return nil, fmt.Errorf("getServicesList error: %w", err)
	}
	return services, nil
}
