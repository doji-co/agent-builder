package generator

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/doji-co/agent-builder/internal/model"
)

func TestGenerator_GenerateProjectFiles(t *testing.T) {
	orchestrator := model.NewOrchestrator("TripCoordinator", model.PatternSequential, "Coordinates travel planning", "gemini-2.5-flash")
	orchestrator.AddSubAgent(model.NewAgent("Researcher", model.AgentTypeLLM, "Research the destination", "research_notes", "gemini-2.5-flash"))
	orchestrator.AddSubAgent(model.NewAgent("ItineraryBuilder", model.AgentTypeCustom, "", "itinerary", ""))

	project := model.NewProject("travel-planner", orchestrator)

	gen := NewGenerator()
	files, err := gen.GenerateProjectFiles(project)
	if err != nil {
		t.Fatalf("GenerateProjectFiles() error = %v", err)
	}

	expectedPaths := []string{
		filepath.Join("travel_planner", "__init__.py"),
		filepath.Join("travel_planner", "agent.py"),
		filepath.Join("travel_planner", ".env.example"),
		filepath.Join("travel_planner", "sub_agents", "researcher", "__init__.py"),
		filepath.Join("travel_planner", "sub_agents", "researcher", "agent.py"),
		filepath.Join("travel_planner", "sub_agents", "itinerary_builder", "__init__.py"),
		filepath.Join("travel_planner", "sub_agents", "itinerary_builder", "agent.py"),
		"requirements.txt",
		"README.md",
	}

	for _, expectedPath := range expectedPaths {
		if _, ok := findFile(files, expectedPath); !ok {
			t.Fatalf("GenerateProjectFiles() missing %s", expectedPath)
		}
	}

	rootInit, _ := findFile(files, filepath.Join("travel_planner", "__init__.py"))
	if strings.TrimSpace(rootInit.Content) != "from . import agent" {
		t.Fatalf("root __init__.py = %q, want %q", strings.TrimSpace(rootInit.Content), "from . import agent")
	}

	rootAgent, _ := findFile(files, filepath.Join("travel_planner", "agent.py"))
	expectedRootStrings := []string{
		"from google.adk.agents import SequentialAgent",
		"from .sub_agents.researcher import agent as researcher",
		"from .sub_agents.itinerary_builder import agent as itinerary_builder",
		"root_agent = SequentialAgent(",
		`name="trip_coordinator"`,
		"sub_agents=[researcher, itinerary_builder]",
	}
	for _, expected := range expectedRootStrings {
		if !strings.Contains(rootAgent.Content, expected) {
			t.Errorf("root agent missing %q", expected)
		}
	}

	readme, _ := findFile(files, "README.md")
	unexpectedReadmeStrings := []string{
		"python main.py",
		"root_agent.run(",
	}
	for _, unexpected := range unexpectedReadmeStrings {
		if strings.Contains(readme.Content, unexpected) {
			t.Errorf("README.md unexpectedly contains %q", unexpected)
		}
	}
	expectedReadmeStrings := []string{
		"adk run travel_planner",
		"adk web --port 8000",
		"travel_planner/.env.example",
	}
	for _, expected := range expectedReadmeStrings {
		if !strings.Contains(readme.Content, expected) {
			t.Errorf("README.md missing %q", expected)
		}
	}

	requirements, _ := findFile(files, "requirements.txt")
	if strings.TrimSpace(requirements.Content) != "google-adk==1.32.0" {
		t.Fatalf("requirements.txt = %q, want %q", strings.TrimSpace(requirements.Content), "google-adk==1.32.0")
	}

	envExample, _ := findFile(files, filepath.Join("travel_planner", ".env.example"))
	expectedEnvStrings := []string{
		"GOOGLE_API_KEY",
		"GOOGLE_GENAI_USE_VERTEXAI=TRUE",
		"GOOGLE_CLOUD_PROJECT",
		"GOOGLE_CLOUD_LOCATION",
	}
	for _, expected := range expectedEnvStrings {
		if !strings.Contains(envExample.Content, expected) {
			t.Errorf(".env.example missing %q", expected)
		}
	}
}

func TestGenerator_GenerateProjectFiles_LLMCoordinator(t *testing.T) {
	orchestrator := model.NewOrchestrator("Coordinator", model.PatternLLMCoordinated, "Delegates requests", "gemini-2.5-flash")
	orchestrator.AddSubAgent(model.NewAgent("Researcher", model.AgentTypeLLM, "Research the topic", "research_notes", "gemini-2.5-flash"))

	project := model.NewProject("assistant", orchestrator)

	gen := NewGenerator()
	files, err := gen.GenerateProjectFiles(project)
	if err != nil {
		t.Fatalf("GenerateProjectFiles() error = %v", err)
	}

	rootAgent, _ := findFile(files, filepath.Join("assistant", "agent.py"))
	expectedStrings := []string{
		"from google.adk.agents import LlmAgent",
		`model="gemini-2.5-flash"`,
		`instruction="""Coordinate the available sub-agents`,
		"sub_agents=[researcher]",
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(rootAgent.Content, expected) {
			t.Errorf("LLM coordinator root agent missing %q", expected)
		}
	}
}

func TestGenerator_GenerateProjectFiles_ParallelPattern(t *testing.T) {
	orchestrator := model.NewOrchestrator("ParallelCoordinator", model.PatternParallel, "Runs tasks in parallel", "gemini-2.5-flash")
	orchestrator.AddSubAgent(model.NewAgent("Researcher", model.AgentTypeLLM, "Research the topic", "research_notes", "gemini-2.5-flash"))

	project := model.NewProject("parallel-assistant", orchestrator)

	gen := NewGenerator()
	files, err := gen.GenerateProjectFiles(project)
	if err != nil {
		t.Fatalf("GenerateProjectFiles() error = %v", err)
	}

	rootAgent, _ := findFile(files, filepath.Join("parallel_assistant", "agent.py"))
	if !strings.Contains(rootAgent.Content, "root_agent = ParallelAgent(") {
		t.Fatalf("parallel root agent should use ParallelAgent")
	}
}

func TestGenerator_GenerateProjectFiles_LoopPattern(t *testing.T) {
	orchestrator := model.NewOrchestrator("LoopCoordinator", model.PatternLoop, "Iterates until done", "gemini-2.5-flash")
	orchestrator.AddSubAgent(model.NewAgent("Researcher", model.AgentTypeLLM, "Research the topic", "research_notes", "gemini-2.5-flash"))

	project := model.NewProject("loop-assistant", orchestrator)

	gen := NewGenerator()
	files, err := gen.GenerateProjectFiles(project)
	if err != nil {
		t.Fatalf("GenerateProjectFiles() error = %v", err)
	}

	rootAgent, _ := findFile(files, filepath.Join("loop_assistant", "agent.py"))
	expectedStrings := []string{
		"root_agent = LoopAgent(",
		"max_iterations=3",
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(rootAgent.Content, expected) {
			t.Fatalf("loop root agent missing %q", expected)
		}
	}
}

func TestGenerator_GenerateSingleAgentFiles_LLM(t *testing.T) {
	agent := model.NewAgent("Researcher", model.AgentTypeLLM, "Research the topic", "research_notes", "gemini-2.5-flash")

	gen := NewGenerator()
	files, err := gen.GenerateSingleAgentFiles(agent)
	if err != nil {
		t.Fatalf("GenerateSingleAgentFiles() error = %v", err)
	}

	expectedPaths := []string{
		filepath.Join("researcher", "__init__.py"),
		filepath.Join("researcher", "agent.py"),
	}
	for _, expectedPath := range expectedPaths {
		if _, ok := findFile(files, expectedPath); !ok {
			t.Fatalf("GenerateSingleAgentFiles() missing %s", expectedPath)
		}
	}

	agentFile, _ := findFile(files, filepath.Join("researcher", "agent.py"))
	expectedStrings := []string{
		"from google.adk.agents import LlmAgent",
		"agent = LlmAgent(",
		`name="researcher"`,
		`model="gemini-2.5-flash"`,
		`output_key="research_notes"`,
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(agentFile.Content, expected) {
			t.Errorf("LLM single agent missing %q", expected)
		}
	}
}

func TestGenerator_GenerateSingleAgentFiles_Custom(t *testing.T) {
	agent := model.NewAgent("Formatter", model.AgentTypeCustom, "", "formatted_result", "")

	gen := NewGenerator()
	files, err := gen.GenerateSingleAgentFiles(agent)
	if err != nil {
		t.Fatalf("GenerateSingleAgentFiles() error = %v", err)
	}

	agentFile, _ := findFile(files, filepath.Join("formatter", "agent.py"))
	expectedStrings := []string{
		"from google.adk.agents import BaseAgent",
		"from google.adk.agents.invocation_context import InvocationContext",
		"from google.adk.events import Event",
		"class FormatterAgent(BaseAgent):",
		"async def _run_async_impl(self, ctx: InvocationContext)",
		`ctx.session.state["formatted_result"] =`,
		"agent = FormatterAgent(",
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(agentFile.Content, expected) {
			t.Errorf("custom single agent missing %q", expected)
		}
	}

	if strings.Contains(agentFile.Content, "LlmAgent(") {
		t.Fatalf("custom single agent should not use LlmAgent")
	}
}

func TestGenerator_GenerateProjectFiles_WithHyphenatedNames(t *testing.T) {
	orchestrator := model.NewOrchestrator("API-Coordinator", model.PatternParallel, "Runs API tasks", "gemini-2.5-flash")
	orchestrator.AddSubAgent(model.NewAgent("grafana-agent", model.AgentTypeLLM, "Query Grafana", "grafana_data", "gemini-2.5-flash"))

	project := model.NewProject("api-project", orchestrator)

	gen := NewGenerator()
	files, err := gen.GenerateProjectFiles(project)
	if err != nil {
		t.Fatalf("GenerateProjectFiles() error = %v", err)
	}

	rootAgent, _ := findFile(files, filepath.Join("api_project", "agent.py"))
	if strings.Contains(rootAgent.Content, "grafana-agent") {
		t.Fatalf("root agent should normalize hyphenated identifiers")
	}
	if !strings.Contains(rootAgent.Content, "from .sub_agents.grafana_agent import agent as grafana_agent") {
		t.Fatalf("root agent should use normalized sub-agent import")
	}
	if !strings.Contains(rootAgent.Content, `name="api_coordinator"`) {
		t.Fatalf("root agent should normalize orchestrator name")
	}
}

func findFile(files []File, path string) (File, bool) {
	for _, file := range files {
		if filepath.Clean(file.Path) == filepath.Clean(path) {
			return file, true
		}
	}
	return File{}, false
}
