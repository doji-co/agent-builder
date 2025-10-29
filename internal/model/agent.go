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
	Tools       []AdkTool
}

func NewAgent(name string, agentType AgentType, instruction, outputKey, model string, adkTools []AdkTool) *Agent {
	return &Agent{
		Name:        name,
		Type:        agentType,
		Instruction: instruction,
		OutputKey:   outputKey,
		Model:       model,
		Tools:       adkTools,
	}
}

func (a *Agent) Validate() error {
	if a.Name == "" {
		return errors.New("name cannot be empty")
	}

	if a.Type == AgentTypeLLM && a.Instruction == "" {
		return errors.New("instruction is required for LLM agents")
	}

	return nil
}
