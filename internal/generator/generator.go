package generator

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/doji-co/agent-builder/internal/model"
)

//go:embed templates/*
var templatesFS embed.FS

type Generator struct {
	templates *template.Template
}

func NewGenerator() *Generator {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"lower":           strings.ToLower,
		"snakeCase":       toSnakeCase,
		"getAgentClass":   getAgentClass,
		"getImports":      getImports,
		"getToolImports":  getToolImports,
		"hasTools":        hasTools,
		"getToolsList":    getToolsList,
	}).ParseFS(templatesFS, "templates/*.tmpl"))

	return &Generator{
		templates: tmpl,
	}
}

func (g *Generator) GenerateAgentPy(project *model.Project) (string, error) {
	var buf bytes.Buffer
	err := g.templates.ExecuteTemplate(&buf, "agent.py.tmpl", project)
	if err != nil {
		return "", fmt.Errorf("failed to generate agent.py: %w", err)
	}
	return buf.String(), nil
}

func (g *Generator) GenerateOrchestratorPy(orchestrator *model.Orchestrator) (string, error) {
	var buf bytes.Buffer
	err := g.templates.ExecuteTemplate(&buf, "orchestrator_agent.py.tmpl", orchestrator)
	if err != nil {
		return "", fmt.Errorf("failed to generate orchestrator agent.py: %w", err)
	}
	return buf.String(), nil
}

func (g *Generator) GenerateSubAgentPy(agent *model.Agent) (string, error) {
	var buf bytes.Buffer
	err := g.templates.ExecuteTemplate(&buf, "agent_single.py.tmpl", agent)
	if err != nil {
		return "", fmt.Errorf("failed to generate sub-agent agent.py: %w", err)
	}
	return buf.String(), nil
}

func (g *Generator) GenerateMainPy(project *model.Project) (string, error) {
	var buf bytes.Buffer
	err := g.templates.ExecuteTemplate(&buf, "main.py.tmpl", project)
	if err != nil {
		return "", fmt.Errorf("failed to generate main.py: %w", err)
	}
	return buf.String(), nil
}

func (g *Generator) GenerateRequirementsTxt() (string, error) {
	var buf bytes.Buffer
	err := g.templates.ExecuteTemplate(&buf, "requirements.txt.tmpl", nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate requirements.txt: %w", err)
	}
	return buf.String(), nil
}

func (g *Generator) GenerateReadme(project *model.Project) (string, error) {
	var buf bytes.Buffer
	err := g.templates.ExecuteTemplate(&buf, "README.md.tmpl", project)
	if err != nil {
		return "", fmt.Errorf("failed to generate README.md: %w", err)
	}
	return buf.String(), nil
}

func toSnakeCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")

	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func getAgentClass(pattern model.OrchestrationPattern) string {
	switch pattern {
	case model.PatternSequential:
		return "SequentialAgent"
	case model.PatternParallel:
		return "ParallelAgent"
	case model.PatternLLMCoordinated:
		return "LlmAgent"
	case model.PatternLoop:
		return "LoopAgent"
	default:
		return "SequentialAgent"
	}
}

func getImports(project *model.Project) string {
	imports := []string{"LlmAgent"}

	agentClass := getAgentClass(project.Orchestrator.Pattern)
	if agentClass != "LlmAgent" {
		imports = append(imports, agentClass)
	}

	return strings.Join(imports, ", ")
}

func getToolImports(agent *model.Agent) string {
	if len(agent.Tools) == 0 {
		return ""
	}

	toolSet := make(map[string]bool)
	for _, tool := range agent.Tools {
		toolSet[string(tool)] = true
	}

	tools := make([]string, 0, len(toolSet))
	for tool := range toolSet {
		tools = append(tools, tool)
	}

	return strings.Join(tools, ", ")
}

func hasTools(agent *model.Agent) bool {
	return len(agent.Tools) > 0
}

func getToolsList(agent *model.Agent) string {
	if len(agent.Tools) == 0 {
		return ""
	}

	tools := make([]string, len(agent.Tools))
	for i, tool := range agent.Tools {
		tools[i] = string(tool)
	}

	return strings.Join(tools, ", ")
}
