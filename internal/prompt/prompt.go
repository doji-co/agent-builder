package prompt

import (
	"errors"
	"regexp"

	"github.com/doji-co/agent-builder/internal/model"
)

const (
	DefaultModel      = "gemini-flash-latest"
	CustomModelOption = "Custom model..."
)

var AvailableModels = []string{
	DefaultModel,
	"gemini-2.5-flash",
	"gemini-2.5-pro",
	CustomModelOption,
}

var (
	projectNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	agentNameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	outputKeyRegex   = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
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

func ValidateOutputKey(key string) error {
	if key == "" {
		return errors.New("output key cannot be empty")
	}
	if !outputKeyRegex.MatchString(key) {
		return errors.New("output key must contain only letters, numbers, and underscores")
	}
	return nil
}

func ResolveModel(selection, custom string) string {
	if selection == CustomModelOption {
		if custom == "" {
			return DefaultModel
		}
		return custom
	}
	if selection == "" {
		return DefaultModel
	}
	return selection
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
