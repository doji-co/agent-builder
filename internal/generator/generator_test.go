package generator

import (
	"strings"
	"testing"

	"github.com/doji-co/agent-builder/internal/model"
)

func TestGenerator_GenerateAgentsPy(t *testing.T) {
	orch := model.NewOrchestrator("ResearchCoordinator", model.PatternSequential, "Coordinates research tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Researcher", model.AgentTypeLLM, "Research the topic", "research_data", "gemini-2.0-flash"))
	orch.AddSubAgent(model.NewAgent("Writer", model.AgentTypeLLM, "Write based on {research_data}", "draft", "gemini-2.0-flash"))

	project := model.NewProject("test-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateAgentsPy(project)

	if err != nil {
		t.Fatalf("GenerateAgentsPy() error = %v", err)
	}

	if content == "" {
		t.Error("GenerateAgentsPy() returned empty content")
	}

	t.Logf("Generated content:\n%s", content)

	expectedStrings := []string{
		"from google.adk.agents import",
		"SequentialAgent",
		"LlmAgent",
		"researcher = LlmAgent(",
		`name="Researcher"`,
		`instruction="Research the topic"`,
		`output_key="research_data"`,
		"writer = LlmAgent(",
		`name="Writer"`,
		"research_coordinator = SequentialAgent(",
		`name="ResearchCoordinator"`,
		"sub_agents=[researcher, writer]",
		"root_agent = research_coordinator",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(content, expected) {
			t.Errorf("GenerateAgentsPy() missing expected string: %s", expected)
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
		"from agents import root_agent",
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

func TestGenerator_GenerateAgentsPy_ParallelPattern(t *testing.T) {
	orch := model.NewOrchestrator("ParallelCoord", model.PatternParallel, "Parallel tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Task1", model.AgentTypeLLM, "Do task 1", "result1", "gemini-2.0-flash"))
	orch.AddSubAgent(model.NewAgent("Task2", model.AgentTypeLLM, "Do task 2", "result2", "gemini-2.0-flash"))

	project := model.NewProject("parallel-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateAgentsPy(project)

	if err != nil {
		t.Fatalf("GenerateAgentsPy() error = %v", err)
	}

	if !strings.Contains(content, "ParallelAgent") {
		t.Error("GenerateAgentsPy() should use ParallelAgent for parallel pattern")
	}

	if !strings.Contains(content, "from google.adk.agents import") && !strings.Contains(content, "ParallelAgent") {
		t.Error("GenerateAgentsPy() missing ParallelAgent import")
	}
}

func TestGenerator_GenerateAgentsPy_LoopPattern(t *testing.T) {
	orch := model.NewOrchestrator("LoopCoord", model.PatternLoop, "Loop tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("Task", model.AgentTypeLLM, "Iterative task", "result", "gemini-2.0-flash"))

	project := model.NewProject("loop-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateAgentsPy(project)

	if err != nil {
		t.Fatalf("GenerateAgentsPy() error = %v", err)
	}

	if !strings.Contains(content, "LoopAgent") {
		t.Error("GenerateAgentsPy() should use LoopAgent for loop pattern")
	}
}

func TestGenerator_GenerateAgentsPy_WithHyphens(t *testing.T) {
	orch := model.NewOrchestrator("APICoordinator", model.PatternSequential, "Coordinates API tasks", "gemini-2.0-flash")
	orch.AddSubAgent(model.NewAgent("grafana-agent", model.AgentTypeLLM, "Query Grafana", "grafana_data", "gemini-2.0-flash"))
	orch.AddSubAgent(model.NewAgent("data-processor", model.AgentTypeLLM, "Process data", "processed", "gemini-2.0-flash"))

	project := model.NewProject("api-project", orch)

	gen := NewGenerator()
	content, err := gen.GenerateAgentsPy(project)

	if err != nil {
		t.Fatalf("GenerateAgentsPy() error = %v", err)
	}

	t.Logf("Generated content:\n%s", content)

	if strings.Contains(content, "grafana-agent = ") {
		t.Error("GenerateAgentsPy() should not create Python variables with hyphens")
	}

	if !strings.Contains(content, "grafana_agent = ") {
		t.Error("GenerateAgentsPy() should convert hyphens to underscores in variable names")
	}

	if !strings.Contains(content, "data_processor = ") {
		t.Error("GenerateAgentsPy() should convert hyphens to underscores in variable names")
	}

	if !strings.Contains(content, `name="grafana-agent"`) {
		t.Error("GenerateAgentsPy() should keep original name in the name field")
	}
}
