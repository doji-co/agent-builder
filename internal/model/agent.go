package model

import "errors"

type AgentType string

const (
	AgentTypeLLM    AgentType = "llm"
	AgentTypeCustom AgentType = "custom"
)

type Agent struct {
	Name        string
	Type        AgentType
	Instruction string
	OutputKey   string
	Model       string
}

func NewAgent(name string, agentType AgentType, instruction, outputKey, model string) *Agent {
	return &Agent{
		Name:        name,
		Type:        agentType,
		Instruction: instruction,
		OutputKey:   outputKey,
		Model:       model,
	}
}

func (a *Agent) Validate() error {
	if a.Name == "" {
		return errors.New("name cannot be empty")
	}

	if a.OutputKey == "" {
		return errors.New("output key cannot be empty")
	}

	if a.Type == AgentTypeLLM && a.Instruction == "" {
		return errors.New("instruction is required for LLM agents")
	}

	if a.Type == AgentTypeLLM && a.Model == "" {
		return errors.New("model is required for LLM agents")
	}

	return nil
}

func (a *Agent) PackageName() string {
	return normalizeName(a.Name)
}

func (a *Agent) ClassName() string {
	return className(a.Name)
}
