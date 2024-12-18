package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"path"
	"strings"
)

type jsonPackage struct {
	ImportPath string
	Name       string
	Deps       []string
}

func getServicesList(ctx context.Context, workingDirectory, goModuleName string) (map[string]string, error) {
	cmd, err := commandContext(ctx, workingDirectory, "go", "list", "-json", "./...")
	if err != nil {
		return nil, fmt.Errorf("commandContext error: %w", err)
	}

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stdout error: %w", err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stderr error: %w", err)
	}

	if err := cmd.Start(); err != nil {
		drainClose(stdOut, stdErr)
		return nil, fmt.Errorf("cmd.Start: run go list: %w", err)
	}

	services := make(map[string]string)
	dec := json.NewDecoder(stdOut)
	for dec.More() {
		var pkg jsonPackage

		if err := dec.Decode(&pkg); err != nil {
			drainClose(stdOut, stdErr)
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
func getReverseDependencies(ctx context.Context, workingDirectory string, onlyServices bool) (map[string][]string, error) {
	cmd, err := commandContext(ctx, workingDirectory, "go", "list", "-json", "./...")
	if err != nil {
		return nil, fmt.Errorf("commandContext error: %w", err)
	}

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stdout: %w", err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("get go list stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		drainClose(stdOut, stdErr)
		return nil, fmt.Errorf("cmd.Start: run go list: %w", err)
	}

	revDeps := map[string][]string{}
	dec := json.NewDecoder(stdOut)
	for dec.More() {
		var pkg jsonPackage
		err = dec.Decode(&pkg)
		if err != nil {
			drainClose(stdOut, stdErr)
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

func drainClose(stdOut, stdErr io.ReadCloser) {
	dataOut, _ := io.ReadAll(stdOut)
	_ = stdOut.Close()
	if len(dataOut) > 0 {
		slog.Warn("stdout", "output", string(dataOut))
	}

	dataErr, _ := io.ReadAll(stdErr)
	_ = stdErr.Close()
	if len(dataErr) > 0 {
		slog.Warn("stderr", "output", string(dataErr))
	}
}
