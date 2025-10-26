package model

import (
	"errors"
	"fmt"
)

type OrchestrationPattern string

const (
	PatternSequential     OrchestrationPattern = "sequential"
	PatternParallel       OrchestrationPattern = "parallel"
	PatternLLMCoordinated OrchestrationPattern = "llm-coordinated"
	PatternLoop           OrchestrationPattern = "loop"
)

func (p OrchestrationPattern) String() string {
	switch p {
	case PatternSequential:
		return "Sequential"
	case PatternParallel:
		return "Parallel"
	case PatternLLMCoordinated:
		return "LLM-Coordinated"
	case PatternLoop:
		return "Loop"
	default:
		return string(p)
	}
}

func (p OrchestrationPattern) Description() string {
	switch p {
	case PatternSequential:
		return "Sub-agents run one after another"
	case PatternParallel:
		return "Sub-agents run simultaneously"
	case PatternLLMCoordinated:
		return "Orchestrator decides which sub-agent to call"
	case PatternLoop:
		return "Repeat sub-agents until condition met"
	default:
		return ""
	}
}

type Orchestrator struct {
	Name        string
	Pattern     OrchestrationPattern
	Description string
	Model       string
	SubAgents   []*Agent
}

func NewOrchestrator(name string, pattern OrchestrationPattern, description, model string) *Orchestrator {
	return &Orchestrator{
		Name:        name,
		Pattern:     pattern,
		Description: description,
		Model:       model,
		SubAgents:   []*Agent{},
	}
}

func (o *Orchestrator) AddSubAgent(agent *Agent) {
	o.SubAgents = append(o.SubAgents, agent)
}

func (o *Orchestrator) Validate() error {
	if o.Name == "" {
		return errors.New("orchestrator name cannot be empty")
	}

	if len(o.SubAgents) == 0 {
		return errors.New("orchestrator must have at least one sub-agent")
	}

	for _, agent := range o.SubAgents {
		if err := agent.Validate(); err != nil {
			return fmt.Errorf("sub-agent validation failed: %w", err)
		}
	}

	return nil
}
