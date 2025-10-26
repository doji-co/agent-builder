package generator

import (
	"strings"
	"testing"

	"github.com/doji-co/agent-builder/internal/model"
)

func TestGenerator_GenerateAgentPy(t *testing.T) {
	orch := model.NewOrchestrator("ResearchCoordinator", model.PatternSequential, "Coordinates research tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Researcher", model.AgentTypeLLM, "Research the topic", "research_data", "gemini-2.0-flash"))
	orch.AddSubAgent(model.NewAgent("Writer", model.AgentTypeLLM, "Write based on {research_data}", "draft", "gemini-2.0-flash"))

	project := model.NewProject("test-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateAgentPy(project)

	if err != nil {
		t.Fatalf("GenerateAgentPy() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateAgentPy() returned empty content")
	}

	t.Logf("Generated content:\n%s", content)

	expectedStrings := []string{
		"from google.adk.agents import",
		"SequentialAgent",
		"LlmAgent",
		"researcher = LlmAgent(",
		`name="researcher"`,
		`instruction="Research the topic"`,
		`output_key="research_data"`,
		"writer = LlmAgent(",
		`name="writer"`,
		"research_coordinator = SequentialAgent(",
		`name="research_coordinator"`,
		"sub_agents=[researcher, writer]",
		"root_agent = research_coordinator",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("GenerateAgentPy() missing expected string: %s", expected)
		}
	}
}

func TestGenerator_GenerateMainPy(t *testing.T) {
	orch := model.NewOrchestrator("Coordinator", model.PatternSequential, "Test", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Agent1", model.AgentTypeLLM, "Task", "result", "gemini-2.0-flash"))

	project := model.NewProject("test-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateMainPy(project)

	if err != nil {
		t.Fatalf("GenerateMainPy() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateMainPy() returned empty content")
	}

	t.Logf("Generated main.py:\n%s", content)

	expectedStrings := []string{
		"from coordinator.agent import agent as root_agent",
		"root_agent.run(",
		`if __name__ == "__main__":`,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("GenerateMainPy() missing expected string: %s", expected)
		}
	}
}

func TestGenerator_GenerateRequirementsTxt(t *testing.T) {
	gen := NewGenerator()
	content, err := gen.GenerateRequirementsTxt()

	if err != nil {
		t.Fatalf("GenerateRequirementsTxt() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateRequirementsTxt() returned empty content")
	}

	if !strings.Contains(content, "google-adk") {
		t.Error("GenerateRequirementsTxt() missing google-adk dependency")
	}
}

func TestGenerator_GenerateReadme(t *testing.T) {
	orch := model.NewOrchestrator("Coordinator", model.PatternSequential, "Test coordinator", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Agent1", model.AgentTypeLLM, "Task", "result", "gemini-2.0-flash"))

	project := model.NewProject("my-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateReadme(project)

	if err != nil {
		t.Fatalf("GenerateReadme() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateReadme() returned empty content")
	}

	expectedStrings := []string{
		"# my-project",
		"pip install -r requirements.txt",
		"python main.py",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("GenerateReadme() missing expected string: %s", expected)
		}
	}
}

func TestGenerator_GenerateAgentPy_ParallelPattern(t *testing.T) {
	orch := model.NewOrchestrator("ParallelCoord", model.PatternParallel, "Parallel tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Task1", model.AgentTypeLLM, "Do task 1", "result1", "gemini-2.0-flash"))
	orch.AddSubAgent(model.NewAgent("Task2", model.AgentTypeLLM, "Do task 2", "result2", "gemini-2.0-flash"))

	project := model.NewProject("parallel-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateAgentPy(project)

	if err != nil {
		t.Fatalf("GenerateAgentPy() error = %v", err)
	}

	if !strings.Contains(content, "ParallelAgent") {
		t.Error("GenerateAgentPy() should use ParallelAgent for parallel pattern")
	}

	if !strings.Contains(content, "from google.adk.agents import") && !strings.Contains(content, "ParallelAgent") {
		t.Error("GenerateAgentPy() missing ParallelAgent import")
	}
}

func TestGenerator_GenerateAgentPy_LoopPattern(t *testing.T) {
	orch := model.NewOrchestrator("LoopCoord", model.PatternLoop, "Loop tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Task", model.AgentTypeLLM, "Iterative task", "result", "gemini-2.0-flash"))

	project := model.NewProject("loop-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateAgentPy(project)

	if err != nil {
		t.Fatalf("GenerateAgentPy() error = %v", err)
	}

	if !strings.Contains(content, "LoopAgent") {
		t.Error("GenerateAgentPy() should use LoopAgent for loop pattern")
	}
}

func TestGenerator_GenerateAgentPy_WithHyphens(t *testing.T) {
	orch := model.NewOrchestrator("APICoordinator", model.PatternSequential, "Coordinates API tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("grafana-agent", model.AgentTypeLLM, "Query Grafana", "grafana_data", "gemini-2.0-flash"))
	orch.AddSubAgent(model.NewAgent("data-processor", model.AgentTypeLLM, "Process data", "processed", "gemini-2.0-flash"))

	project := model.NewProject("api-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateAgentPy(project)

	if err != nil {
		t.Fatalf("GenerateAgentPy() error = %v", err)
	}

	t.Logf("Generated content:\n%s", content)

	if strings.Contains(content, "grafana-agent = ") {
		t.Error("GenerateAgentPy() should not create Python variables with hyphens")
	}

	if !strings.Contains(content, "grafana_agent = ") {
		t.Error("GenerateAgentPy() should convert hyphens to underscores in variable names")
	}

	if !strings.Contains(content, "data_processor = ") {
		t.Error("GenerateAgentPy() should convert hyphens to underscores in variable names")
	}

	if !strings.Contains(content, `name="grafana_agent"`) {
		t.Error("GenerateAgentPy() should convert hyphens to underscores in the name field")
	}
}

func TestGenerator_GenerateOrchestratorPy(t *testing.T) {
	orch := model.NewOrchestrator("ResearchCoordinator", model.PatternSequential, "Coordinates research tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Researcher", model.AgentTypeLLM, "Research the topic", "research_data", "gemini-2.0-flash"))
	orch.AddSubAgent(model.NewAgent("Writer", model.AgentTypeLLM, "Write based on research", "draft", "gemini-2.0-flash"))

	gen := NewGenerator()
	content, err := gen.GenerateOrchestratorPy(orch)

	if err != nil {
		t.Fatalf("GenerateOrchestratorPy() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateOrchestratorPy() returned empty content")
	}

	t.Logf("Generated orchestrator:\n%s", content)

	expectedStrings := []string{
		"from google.adk.agents import SequentialAgent",
		"from researcher.agent import agent as researcher",
		"from writer.agent import agent as writer",
		"agent = SequentialAgent(",
		`name="research_coordinator"`,
		"sub_agents=[researcher, writer]",
		"root_agent = agent",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("GenerateOrchestratorPy() missing expected string: %s", expected)
		}
	}
}

func TestGenerator_GenerateSubAgentPy(t *testing.T) {
	agent := model.NewAgent("Researcher", model.AgentTypeLLM, "Research the topic", "research_data", "gemini-2.0-flash")

	gen := NewGenerator()
	content, err := gen.GenerateSubAgentPy(agent)

	if err != nil {
		t.Fatalf("GenerateSubAgentPy() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateSubAgentPy() returned empty content")
	}

	t.Logf("Generated sub-agent:\n%s", content)

	expectedStrings := []string{
		"from google.adk.agents import LlmAgent",
		"agent = LlmAgent(",
		`name="researcher"`,
		`instruction="Research the topic"`,
		`output_key="research_data"`,
		`model="gemini-2.0-flash"`,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("GenerateSubAgentPy() missing expected string: %s", expected)
		}
	}
}
