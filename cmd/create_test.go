package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/doji-co/agent-builder/internal/model"
)

func TestCreateCommand_GeneratesFullProject(t *testing.T) {
	tempDir := t.TempDir()
	prompter := &stubPrompter{
		projectType:             "full",
		projectName:             "travel-planner",
		orchestrationPattern:    model.PatternSequential,
		orchestratorName:        "TripCoordinator",
		orchestratorDescription: "Coordinates travel planning",
		orchestratorModel:       "gemini-2.5-flash",
		agents: []stubAgentInput{
			{
				name:        "Researcher",
				agentType:   model.AgentTypeLLM,
				instruction: "Research the destination",
				outputKey:   "research_notes",
				model:       "gemini-2.5-flash",
			},
			{
				name:      "Formatter",
				agentType: model.AgentTypeCustom,
				outputKey: "formatted_result",
			},
		},
		outputDirectory: filepath.Join(tempDir, "travel-planner"),
	}

	cmd := newCreateCommand(createDependencies{
		prompterFactory: func() createPrompter { return prompter },
	})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	rootAgentPath := filepath.Join(tempDir, "travel-planner", "travel_planner", "agent.py")
	rootAgentContent, err := os.ReadFile(rootAgentPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", rootAgentPath, err)
	}

	expectedStrings := []string{
		"root_agent = SequentialAgent(",
		"from .sub_agents.researcher import agent as researcher",
		"from .sub_agents.formatter import agent as formatter",
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(string(rootAgentContent), expected) {
			t.Fatalf("root agent content = %q, want to contain %q", string(rootAgentContent), expected)
		}
	}

	if !strings.Contains(buf.String(), "adk run travel_planner") {
		t.Fatalf("create output = %q, want to contain %q", buf.String(), "adk run travel_planner")
	}
}

func TestCreateCommand_GeneratesSingleCustomAgentPackage(t *testing.T) {
	tempDir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})

	prompter := &stubPrompter{
		projectType: "single",
		agents: []stubAgentInput{
			{
				name:      "Formatter",
				agentType: model.AgentTypeCustom,
				outputKey: "formatted_result",
			},
		},
	}

	cmd := newCreateCommand(createDependencies{
		prompterFactory: func() createPrompter { return prompter },
	})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	agentPath := filepath.Join(tempDir, "formatter", "agent.py")
	agentContent, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", agentPath, err)
	}

	if !strings.Contains(string(agentContent), "class FormatterAgent(BaseAgent):") {
		t.Fatalf("single custom agent content = %q, want custom BaseAgent scaffold", string(agentContent))
	}
	if strings.Contains(string(agentContent), "LlmAgent(") {
		t.Fatalf("single custom agent content should not contain LlmAgent")
	}
}

type stubPrompter struct {
	projectType             string
	projectName             string
	orchestrationPattern    model.OrchestrationPattern
	orchestratorName        string
	orchestratorDescription string
	orchestratorModel       string
	agents                  []stubAgentInput
	outputDirectory         string
	addAnotherIndex         int
}

type stubAgentInput struct {
	name        string
	agentType   model.AgentType
	instruction string
	outputKey   string
	model       string
}

func (s *stubPrompter) PromptProjectType() (string, error) {
	return s.projectType, nil
}

func (s *stubPrompter) PromptProjectName() (string, error) {
	return s.projectName, nil
}

func (s *stubPrompter) PromptOrchestrationPattern() (model.OrchestrationPattern, error) {
	return s.orchestrationPattern, nil
}

func (s *stubPrompter) PromptOrchestratorName() (string, error) {
	return s.orchestratorName, nil
}

func (s *stubPrompter) PromptOrchestratorDescription() (string, error) {
	return s.orchestratorDescription, nil
}

func (s *stubPrompter) PromptModel(defaultModel string) (string, error) {
	if s.projectType == "full" && s.orchestratorModel != "" && s.addAnotherIndex == 0 {
		model := s.orchestratorModel
		s.orchestratorModel = ""
		return model, nil
	}
	index := s.addAnotherIndex
	if index >= len(s.agents) {
		index = len(s.agents) - 1
	}
	if index < 0 {
		return defaultModel, nil
	}
	if s.agents[index].model == "" {
		return defaultModel, nil
	}
	return s.agents[index].model, nil
}

func (s *stubPrompter) PromptAgentName(agentNumber int) (string, error) {
	return s.agents[agentNumber-1].name, nil
}

func (s *stubPrompter) PromptAgentType() (model.AgentType, error) {
	return s.agents[s.addAnotherIndex].agentType, nil
}

func (s *stubPrompter) PromptAgentInstruction(agentName string) (string, error) {
	return s.agents[s.addAnotherIndex].instruction, nil
}

func (s *stubPrompter) PromptOutputKey() (string, error) {
	return s.agents[s.addAnotherIndex].outputKey, nil
}

func (s *stubPrompter) PromptAddAnotherAgent() (bool, error) {
	s.addAnotherIndex++
	return s.addAnotherIndex < len(s.agents), nil
}

func (s *stubPrompter) PromptOutputDirectory(defaultDir string) (string, error) {
	return s.outputDirectory, nil
}
