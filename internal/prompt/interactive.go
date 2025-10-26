package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/doji-co/agent-builder/internal/model"
)

type Interactive struct{}

func NewInteractive() *Interactive {
	return &Interactive{}
}

func (i *Interactive) PromptProjectName() (string, error) {
	var name string
	prompt := &survey.Input{
		Message: "Project name?",
	}
	err := survey.AskOne(prompt, &name, survey.WithValidator(func(val interface{}) error {
		if str, ok := val.(string); ok {
			return ValidateProjectName(str)
		}
		return fmt.Errorf("invalid input type")
	}))
	return name, err
}

func (i *Interactive) PromptOrchestrationPattern() (model.OrchestrationPattern, error) {
	patterns := GetOrchestrationPatterns()
	options := make([]string, len(patterns))
	for idx, p := range patterns {
		options[idx] = fmt.Sprintf("%s (%s)", p.String(), p.Description())
	}

	var selection string
	prompt := &survey.Select{
		Message: "Choose orchestration pattern:",
		Options: options,
	}
	err := survey.AskOne(prompt, &selection)
	if err != nil {
		return "", err
	}

	for idx, opt := range options {
		if opt == selection {
			return patterns[idx], nil
		}
	}

	return "", fmt.Errorf("invalid selection")
}

func (i *Interactive) PromptOrchestratorName() (string, error) {
	var name string
	prompt := &survey.Input{
		Message: "Orchestrator name?",
	}
	err := survey.AskOne(prompt, &name, survey.WithValidator(func(val interface{}) error {
		if str, ok := val.(string); ok {
			return ValidateAgentName(str)
		}
		return fmt.Errorf("invalid input type")
	}))
	return name, err
}

func (i *Interactive) PromptOrchestratorDescription() (string, error) {
	var description string
	prompt := &survey.Input{
		Message: "Orchestrator description?",
	}
	err := survey.AskOne(prompt, &description)
	return description, err
}

func (i *Interactive) PromptModel(defaultModel string) (string, error) {
	var modelName string
	prompt := &survey.Input{
		Message: "Model?",
		Default: defaultModel,
	}
	err := survey.AskOne(prompt, &modelName)
	if modelName == "" {
		modelName = defaultModel
	}
	return modelName, err
}

func (i *Interactive) PromptAgentName(agentNumber int) (string, error) {
	var name string
	prompt := &survey.Input{
		Message: fmt.Sprintf("Sub-agent #%d name?", agentNumber),
	}
	err := survey.AskOne(prompt, &name, survey.WithValidator(func(val interface{}) error {
		if str, ok := val.(string); ok {
			return ValidateAgentName(str)
		}
		return fmt.Errorf("invalid input type")
	}))
	return name, err
}

func (i *Interactive) PromptAgentType() (model.AgentType, error) {
	types := GetAgentTypes()
	options := []string{
		"LLM Agent (powered by language model)",
		"Custom Agent (your own Python class)",
	}

	var selection string
	prompt := &survey.Select{
		Message: "Agent type:",
		Options: options,
	}
	err := survey.AskOne(prompt, &selection)
	if err != nil {
		return "", err
	}

	if selection == options[0] {
		return types[0], nil
	}
	return types[1], nil
}

func (i *Interactive) PromptAgentInstruction(agentName string) (string, error) {
	var instruction string
	prompt := &survey.Input{
		Message: fmt.Sprintf("Instruction for %s?", agentName),
	}
	err := survey.AskOne(prompt, &instruction)
	return instruction, err
}

func (i *Interactive) PromptOutputKey() (string, error) {
	var key string
	prompt := &survey.Input{
		Message: "Output key? (where to store result)",
	}
	err := survey.AskOne(prompt, &key)
	return key, err
}

func (i *Interactive) PromptAddAnotherAgent() (bool, error) {
	var add bool
	prompt := &survey.Confirm{
		Message: "Add another sub-agent?",
		Default: true,
	}
	err := survey.AskOne(prompt, &add)
	return add, err
}

func (i *Interactive) PromptOutputDirectory(defaultDir string) (string, error) {
	var dir string
	prompt := &survey.Input{
		Message: "Output directory?",
		Default: defaultDir,
	}
	err := survey.AskOne(prompt, &dir)
	if dir == "" {
		dir = defaultDir
	}
	return dir, err
}

func (i *Interactive) PromptAddExample() (bool, error) {
	var add bool
	prompt := &survey.Confirm{
		Message: "Generate example usage?",
		Default: true,
	}
	err := survey.AskOne(prompt, &add)
	return add, err
}

func (i *Interactive) PromptAddDocker() (bool, error) {
	var add bool
	prompt := &survey.Confirm{
		Message: "Add Docker support?",
		Default: false,
	}
	err := survey.AskOne(prompt, &add)
	return add, err
}
