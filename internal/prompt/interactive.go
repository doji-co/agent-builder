package prompt

import (
	"fmt"

	"charm.land/huh/v2"
	"github.com/doji-co/agent-builder/internal/model"
)

type Interactive struct{}

func NewInteractive() *Interactive {
	return &Interactive{}
}

func (i *Interactive) PromptProjectType() (string, error) {
	var selection string
	err := huh.NewSelect[string]().
		Title("What would you like to create?").
		Options(
			huh.NewOption("Starter project (orchestrator + sub-agents)", "full"),
			huh.NewOption("Single agent (add to existing project)", "single"),
		).
		Value(&selection).
		Run()
	return selection, err
}

func (i *Interactive) PromptProjectName() (string, error) {
	var name string
	err := huh.NewInput().
		Title("Project name").
		Placeholder("travel-planner").
		Validate(ValidateProjectName).
		Value(&name).
		Run()
	return name, err
}

func (i *Interactive) PromptOrchestrationPattern() (model.OrchestrationPattern, error) {
	var selection model.OrchestrationPattern
	err := huh.NewSelect[model.OrchestrationPattern]().
		Title("Choose orchestration pattern").
		Options(
			huh.NewOption(fmt.Sprintf("%s (%s)", model.PatternSequential.String(), model.PatternSequential.Description()), model.PatternSequential),
			huh.NewOption(fmt.Sprintf("%s (%s)", model.PatternParallel.String(), model.PatternParallel.Description()), model.PatternParallel),
			huh.NewOption(fmt.Sprintf("%s (%s)", model.PatternLLMCoordinated.String(), model.PatternLLMCoordinated.Description()), model.PatternLLMCoordinated),
			huh.NewOption(fmt.Sprintf("%s (%s)", model.PatternLoop.String(), model.PatternLoop.Description()), model.PatternLoop),
		).
		Value(&selection).
		Run()
	return selection, err
}

func (i *Interactive) PromptOrchestratorName() (string, error) {
	var name string
	err := huh.NewInput().
		Title("Orchestrator name").
		Placeholder("TripCoordinator").
		Validate(ValidateAgentName).
		Value(&name).
		Run()
	return name, err
}

func (i *Interactive) PromptOrchestratorDescription() (string, error) {
	var description string
	err := huh.NewText().
		Title("Orchestrator description").
		Value(&description).
		Run()
	return description, err
}

func (i *Interactive) PromptModel(defaultModel string) (string, error) {
	var selection string
	err := huh.NewSelect[string]().
		Title("Choose model").
		Options(
			huh.NewOption(AvailableModels[0], AvailableModels[0]),
			huh.NewOption(AvailableModels[1], AvailableModels[1]),
			huh.NewOption(AvailableModels[2], AvailableModels[2]),
			huh.NewOption(AvailableModels[3], AvailableModels[3]),
		).
		Value(&selection).
		Run()
	if err != nil {
		return "", err
	}

	if selection != CustomModelOption {
		return ResolveModel(selection, ""), nil
	}

	var custom string
	err = huh.NewInput().
		Title("Custom model").
		Placeholder(defaultModel).
		Value(&custom).
		Run()
	if err != nil {
		return "", err
	}

	return ResolveModel(selection, custom), nil
}

func (i *Interactive) PromptAgentName(agentNumber int) (string, error) {
	var name string
	err := huh.NewInput().
		Title(fmt.Sprintf("Sub-agent #%d name", agentNumber)).
		Placeholder("Researcher").
		Validate(ValidateAgentName).
		Value(&name).
		Run()
	return name, err
}

func (i *Interactive) PromptAgentType() (model.AgentType, error) {
	var selection model.AgentType
	err := huh.NewSelect[model.AgentType]().
		Title("Agent type").
		Options(
			huh.NewOption("LLM Agent", model.AgentTypeLLM),
			huh.NewOption("Custom Agent", model.AgentTypeCustom),
		).
		Value(&selection).
		Run()
	return selection, err
}

func (i *Interactive) PromptAgentInstruction(agentName string) (string, error) {
	var instruction string
	err := huh.NewText().
		Title(fmt.Sprintf("Instruction for %s", agentName)).
		Value(&instruction).
		Run()
	return instruction, err
}

func (i *Interactive) PromptOutputKey() (string, error) {
	var key string
	err := huh.NewInput().
		Title("Output key").
		Placeholder("research_notes").
		Validate(ValidateOutputKey).
		Value(&key).
		Run()
	return key, err
}

func (i *Interactive) PromptAddAnotherAgent() (bool, error) {
	var add bool
	err := huh.NewConfirm().
		Title("Add another sub-agent?").
		Value(&add).
		Run()
	return add, err
}

func (i *Interactive) PromptOutputDirectory(defaultDir string) (string, error) {
	var dir string
	err := huh.NewInput().
		Title("Output directory").
		Placeholder(defaultDir).
		Value(&dir).
		Run()
	if dir == "" {
		dir = defaultDir
	}
	return dir, err
}
