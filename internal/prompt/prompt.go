package prompt

import (
	"errors"
	"regexp"

	"github.com/doji-co/agent-builder/internal/model"
)

const DefaultModel = "gemini-2.5-flash"

var AvailableModels = []string{
	"gemini-2.5-flash",
	"gemini-2.5-pro",
	"gemini-2.5-flash-lite",
}

var (
	projectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	agentNameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

func ValidateProjectName(name string) error {
	if name == "" {
		return errors.New("project name cannot be empty")
	}
	if !projectNameRegex.MatchString(name) {
		return errors.New("project name must contain only letters, numbers, hyphens, and underscores")
	}
	return nil
}

func ValidateAgentName(name string) error {
	if name == "" {
		return errors.New("agent name cannot be empty")
	}
	if !agentNameRegex.MatchString(name) {
		return errors.New("agent name must contain only letters, numbers, hyphens, and underscores")
	}
	return nil
}

func GetOrchestrationPatterns() []model.OrchestrationPattern {
	return []model.OrchestrationPattern{
		model.PatternSequential,
		model.PatternParallel,
		model.PatternLLMCoordinated,
		model.PatternLoop,
	}
}

func GetAgentTypes() []model.AgentType {
	return []model.AgentType{
		model.AgentTypeLLM,
		model.AgentTypeCustom,
	}
}
