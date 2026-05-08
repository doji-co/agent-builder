package model

import (
	"errors"
	"fmt"
)

type Project struct {
	Name         string
	Orchestrator *Orchestrator
	OutputDir    string
	AddReadme    bool
}

func NewProject(name string, orchestrator *Orchestrator) *Project {
	return &Project{
		Name:         name,
		Orchestrator: orchestrator,
		OutputDir:    fmt.Sprintf("./%s", name),
		AddReadme:    true,
	}
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return errors.New("project name cannot be empty")
	}

	if p.Orchestrator == nil {
		return errors.New("orchestrator cannot be nil")
	}

	if err := p.Orchestrator.Validate(); err != nil {
		return fmt.Errorf("orchestrator validation failed: %w", err)
	}

	return nil
}

func (p *Project) PackageName() string {
	return normalizeName(p.Name)
}
