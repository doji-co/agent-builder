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

func (i *Interactive) PromptProjectType() (string, error) {
	fmt.Println("\n💡 What would you like to create?")
	fmt.Println("   • Starter Project: Complete multi-agent system with orchestrator and sub-agents")
	fmt.Println("     Use this when starting a new ADK project from scratch")
	fmt.Println()
	fmt.Println("   • Single Agent: Just one agent folder to add to an existing project")
	fmt.Println("     Use this when you want to add a new sub-agent to a project you already have")
	fmt.Println()

	var selection string
	prompt := &survey.Select{
		Message: "What would you like to create?",
		Options: []string{
			"Starter project (orchestrator + sub-agents)",
			"Single agent (add to existing project)",
		},
		Help: "Choose based on whether you're starting fresh or extending an existing project",
	}
	err := survey.AskOne(prompt, &selection)
	if err != nil {
		return "", err
	}

	if selection == "Starter project (orchestrator + sub-agents)" {
		return "full", nil
	}
	return "single", nil
}

func (i *Interactive) PromptProjectName() (string, error) {
	fmt.Println("\n💡 What is a project?")
	fmt.Println("   A project is a complete multi-agent system. It will contain all your agents")
	fmt.Println("   and their configuration. Use kebab-case (my-project) or snake_case (my_project).")
	fmt.Println()

	var name string
	prompt := &survey.Input{
		Message: "Project name?",
		Help:    "Example: research-assistant, data-processor, content-generator",
	}
	err := survey.AskOne(prompt, &name, survey.WithValidator(func(val any) error {
		if str, ok := val.(string); ok {
			return ValidateProjectName(str)
		}
		return fmt.Errorf("invalid input type")
	}))
	return name, err
}

func (i *Interactive) PromptOrchestrationPattern() (model.OrchestrationPattern, error) {
	fmt.Println("\n💡 What is an orchestration pattern?")
	fmt.Println("   The pattern determines HOW your agents work together:")
	fmt.Println("   • Sequential: Agents run one after another (like an assembly line)")
	fmt.Println("   • Parallel: Agents run at the same time (for independent tasks)")
	fmt.Println("   • LLM-Coordinated: The orchestrator decides which agent to call")
	fmt.Println("   • Loop: Agents repeat until a condition is met (for refinement)")
	fmt.Println()

	patterns := GetOrchestrationPatterns()
	options := make([]string, len(patterns))
	for idx, p := range patterns {
		options[idx] = fmt.Sprintf("%s (%s)", p.String(), p.Description())
	}

	var selection string
	prompt := &survey.Select{
		Message: "Choose orchestration pattern:",
		Options: options,
		Help:    "Most common: Sequential (for pipelines) or Parallel (for concurrent tasks)",
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
	fmt.Println("\n💡 What is an orchestrator?")
	fmt.Println("   The orchestrator is the ROOT agent that manages all sub-agents.")
	fmt.Println("   It coordinates when and how sub-agents execute their tasks.")
	fmt.Println()
	fmt.Println("   📝 Best practices:")
	fmt.Println("   • Use descriptive names that indicate the system's purpose")
	fmt.Println("   • Common patterns: [Purpose]Coordinator, [Domain]Orchestrator, [Task]Manager")
	fmt.Println("   • Examples: ResearchCoordinator, DataPipelineOrchestrator, ContentManager")
	fmt.Println()

	var name string
	prompt := &survey.Input{
		Message: "Orchestrator name?",
		Help:    "This will be the main agent that controls your system",
	}
	err := survey.AskOne(prompt, &name, survey.WithValidator(func(val any) error {
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
	fmt.Println("\n💡 What is a model?")
	fmt.Println("   The model is the AI that powers the agent's intelligence.")
	fmt.Println()
	fmt.Println("   📊 Available models:")
	fmt.Println("   • gemini-2.5-flash: Fast and efficient (recommended for most use cases)")
	fmt.Println("   • gemini-2.5-pro: Most capable, best for complex reasoning")
	fmt.Println("   • gemini-2.5-flash-lite: Fastest, best for simple tasks")
	fmt.Println()

	var selection string
	prompt := &survey.Select{
		Message: "Choose model:",
		Options: AvailableModels,
		Default: defaultModel,
		Help:    "Start with gemini-2.5-flash and upgrade to pro if needed",
	}
	err := survey.AskOne(prompt, &selection)
	return selection, err
}

func (i *Interactive) PromptAgentName(agentNumber int) (string, error) {
	if agentNumber == 1 {
		fmt.Println("\n💡 What are sub-agents?")
		fmt.Println("   Sub-agents are specialized agents that perform specific tasks.")
		fmt.Println("   The orchestrator coordinates these agents to accomplish complex goals.")
		fmt.Println()
		fmt.Println("   📝 Naming best practices:")
		fmt.Println("   • Use names that describe the agent's specific role")
		fmt.Println("   • Examples: Researcher, Writer, Reviewer, DataFetcher, Analyzer")
		fmt.Println("   • Can use kebab-case (data-processor) or PascalCase (DataProcessor)")
		fmt.Println()
	}

	var name string
	prompt := &survey.Input{
		Message: fmt.Sprintf("Sub-agent #%d name?", agentNumber),
		Help:    "What specific task will this agent perform?",
	}
	err := survey.AskOne(prompt, &name, survey.WithValidator(func(val any) error {
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
	fmt.Println("\n💡 What is an instruction?")
	fmt.Println("   The instruction tells the agent WHAT to do. Be specific and clear.")
	fmt.Println("   The agent will use this as its main goal when processing tasks.")
	fmt.Println()
	fmt.Println("   📝 Examples:")
	fmt.Println("   • 'Research the given topic and provide key findings'")
	fmt.Println("   • 'Write a comprehensive article based on the research data'")
	fmt.Println("   • 'Review the content for quality and suggest improvements'")
	fmt.Println()

	var instruction string
	prompt := &survey.Input{
		Message: fmt.Sprintf("Instruction for %s?", agentName),
		Help:    "Be specific about what this agent should accomplish",
	}
	err := survey.AskOne(prompt, &instruction)
	return instruction, err
}

func (i *Interactive) PromptOutputKey() (string, error) {
	fmt.Println("\n💡 What is an output key?")
	fmt.Println("   The output key is WHERE the agent stores its result for other agents.")
	fmt.Println("   Subsequent agents can reference this data using {output_key} in their instructions.")
	fmt.Println()
	fmt.Println("   📝 Best practices:")
	fmt.Println("   • Use snake_case: research_data, processed_text, final_report")
	fmt.Println("   • Be descriptive: what kind of data does this agent produce?")
	fmt.Println("   • Examples: article_draft, analysis_results, review_feedback")
	fmt.Println()

	var key string
	prompt := &survey.Input{
		Message: "Output key?",
		Help:    "Use snake_case to name where this agent's result will be stored",
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
