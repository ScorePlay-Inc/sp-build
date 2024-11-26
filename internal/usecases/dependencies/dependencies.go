package dependencies

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

func repoRootPath() (string, error) {
	rootPath, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("exec.Command error: %w", err)
	}
	return strings.TrimSpace(string(rootPath)), nil
}

func modifiedFilesSinceLastCommit() ([]string, error) {
	modifiedFiles, err := exec.Command("git", "diff", "--name-only", "@^").Output()
	if err != nil {
		return nil, fmt.Errorf("exec.Command error: %w", err)
	}
	return strings.Split(strings.TrimSpace(string(modifiedFiles)), "\n"), nil
}

func ModifiedPackages(ctx context.Context, onlyServices bool) ([]string, error) {
	repoRoot, err := repoRootPath()
	if err != nil {
		return nil, fmt.Errorf("getRootPath error: %w", err)
	}

	modifiedFiles, err := modifiedFilesSinceLastCommit()
	if err != nil {
		return nil, fmt.Errorf("modifiedFilesSinceLastCommit error %w", err)
	}

	revDeps, err := getReverseDependencies(ctx, repoRoot, onlyServices)
	if err != nil {
		return nil, fmt.Errorf("getReverseDependencies error %w", err)
	}

	toBuild := make(map[string]bool)
	for _, file := range modifiedFiles {
		if strings.HasPrefix(file, "vendor/") {
			file = strings.TrimPrefix(path.Dir(file), "vendor/")
		} else {
			file = path.Join("github.com/ScorePlay-Inc/media-management", path.Dir(file))
		}

		for _, service := range revDeps[file] {
			toBuild[service] = true
		}
	}

	serviceList := make([]string, 0)
	for service := range toBuild {
		prefix := "github.com/ScorePlay-Inc/media-management"
		serviceName := strings.TrimPrefix(service, prefix)
		if serviceName == "tools/revdeps" {
			continue
		}
		serviceList = append(serviceList, serviceName)
	}
	return serviceList, nil
}

func ServicesList(ctx context.Context) (map[string]string, error) {
	repoRoot, err := repoRootPath()
	if err != nil {
		return nil, fmt.Errorf("repoRootPath error: %w", err)
	}
	repoRoot = path.Join(repoRoot, "/app")
	services, err := getServicesList(ctx, repoRoot)
	if err != nil {
		return nil, fmt.Errorf("getServicesList error: %w", err)
	}
	return services, nil
}
