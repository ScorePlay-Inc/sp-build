package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func golangModuleName(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "go", "mod", "edit", "--json")

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
	return module.Module.Path, nil
}

func modifiedFilesSinceLastCommit(ctx context.Context) ([]string, error) {
	modifiedFiles, err := exec.CommandContext(ctx, "git", "diff", "--relative", "--name-only", "@^").Output()
	if err != nil {
		return nil, fmt.Errorf("exec.Command error: %w", err)
	}
	return strings.Split(strings.TrimSpace(string(modifiedFiles)), "\n"), nil
}

func modifiedPackages(ctx context.Context, onlyServices bool) ([]string, error) {
	moduleName, err := golangModuleName(ctx)
	if err != nil {
		return nil, fmt.Errorf("golangModuleName error %w", err)
	}

	modifiedFiles, err := modifiedFilesSinceLastCommit(ctx)
	if err != nil {
		return nil, fmt.Errorf("modifiedFilesSinceLastCommit error %w", err)
	}

	revDeps, err := getReverseDependencies(ctx, onlyServices)
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
		pkgName := strings.TrimPrefix(pkg, moduleName)
		if pkgName == "" {
			continue
		}

		packageList = append(packageList, "."+pkgName)
	}
	return packageList, nil
}

func servicesList(ctx context.Context) (map[string]string, error) {
	goModuleName, err := golangModuleName(ctx)
	if err != nil {
		return nil, fmt.Errorf("golangModuleName error %w", err)
	}

	services, err := getServicesList(ctx, goModuleName)
	if err != nil {
		return nil, fmt.Errorf("getServicesList error: %w", err)
	}
	return services, nil
}
