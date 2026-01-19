package service

import (
	"os/exec"
	"strings"
)

type CommandExecutor interface {
	Execute(name string, args ...string) (string, error)
}

type RealCommandExecutor struct{}

func (r *RealCommandExecutor) Execute(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

type GcloudService struct {
	executor CommandExecutor
}

func NewGcloudService(executor CommandExecutor) *GcloudService {
	return &GcloudService{
		executor: executor,
	}
}

func NewDefaultGcloudService() *GcloudService {
	return NewGcloudService(&RealCommandExecutor{})
}

func (g *GcloudService) IsAvailable() bool {
	_, err := g.executor.Execute("gcloud", "version")
	return err == nil
}

func (g *GcloudService) GetProjectID() (string, error) {
	output, err := g.executor.Execute("gcloud", "config", "get-value", "project")
	if err != nil {
		return "", err
	}
	output = strings.TrimSpace(output)
	if output == "(unset)" {
		return "", nil
	}
	return output, nil
}

func (g *GcloudService) GetRegion() (string, error) {
	output, err := g.executor.Execute("gcloud", "config", "get-value", "compute/region")
	if err != nil {
		return "", err
	}
	output = strings.TrimSpace(output)
	if output == "(unset)" {
		return "", nil
	}
	return output, nil
}
