package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
)

type jsonPackage struct {
	ImportPath string
	Name       string
	Deps       []string
}

func getServicesList(ctx context.Context, goModuleName string) (map[string]string, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cmd := exec.CommandContext(ctx, "go", "list", "-json", "./...")

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stdout error: %w", err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stderr error: %w", err)
	}

	if err := cmd.Start(); err != nil {
		drainClose(logger, stdOut, stdErr)
		return nil, fmt.Errorf("cmd.Start: run go list: %w", err)
	}

	services := make(map[string]string)
	dec := json.NewDecoder(stdOut)
	for dec.More() {
		var pkg jsonPackage
		err = dec.Decode(&pkg)
		if err != nil {
			drainClose(logger, stdOut, stdErr)
			return nil, fmt.Errorf("parse go list output: %w", err)
		}

		if pkg.Name != "main" {
			continue
		}
		servicePath := strings.TrimPrefix(pkg.ImportPath, goModuleName)
		serviceName := path.Base(servicePath)
		services[serviceName] = path.Join(goModuleName, servicePath)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("cmd.Wait:run go list: %w", err)
	}
	return services, nil
}

// getReverseDependencies gets all the executable (package main) dependencies of every package in the repo.
func getReverseDependencies(ctx context.Context, onlyServices bool) (map[string][]string, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cmd := exec.CommandContext(ctx, "go", "list", "-json", "./...")

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stdout: %w", err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		drainClose(logger, stdOut, stdErr)
		return nil, fmt.Errorf("cmd.Start: run go list: %w", err)
	}

	revDeps := map[string][]string{}
	dec := json.NewDecoder(stdOut)
	for dec.More() {
		var pkg jsonPackage
		err = dec.Decode(&pkg)
		if err != nil {
			drainClose(logger, stdOut, stdErr)
			return nil, fmt.Errorf("parse go list output: %w", err)
		}

		if onlyServices && pkg.Name != "main" {
			continue
		}

		revDeps[pkg.ImportPath] = append(revDeps[pkg.ImportPath], pkg.ImportPath)
		for _, dep := range pkg.Deps {
			dep = strings.TrimPrefix(dep, "vendor/")

			if !strings.Contains(strings.Split(dep, "/")[0], ".") {
				continue
			}
			revDeps[dep] = append(revDeps[dep], pkg.ImportPath)
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("cmd.Wait: run go list: %w", err)
	}
	return revDeps, nil
}

func drainClose(logger *slog.Logger, readers ...io.ReadCloser) {
	for _, reader := range readers {
		bts, err := io.ReadAll(reader)
		if err != nil {
			_ = reader.Close()
			continue
		}
		logger.Info(string(bts))
		_ = reader.Close()
	}
}
