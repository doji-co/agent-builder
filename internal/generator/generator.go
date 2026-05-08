package generator

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/doji-co/agent-builder/internal/model"
)

//go:embed templates/*
var templatesFS embed.FS

type File struct {
	Path    string
	Content string
}

type Generator struct {
	templates *template.Template
}

func NewGenerator() *Generator {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"snakeCase":              snakeCase,
		"className":              className,
		"getAgentClass":          getAgentClass,
		"getOrchestratorImports": getOrchestratorImports,
		"coordinatorInstruction": coordinatorInstruction,
	}).ParseFS(templatesFS, "templates/*.tmpl"))

	return &Generator{
		templates: tmpl,
	}
}

func (g *Generator) GenerateProjectFiles(project *model.Project) ([]File, error) {
	rootPackage := project.PackageName()

	rootAgent, err := g.renderTemplate("root_agent.py.tmpl", project)
	if err != nil {
		return nil, err
	}

	envExample, err := g.renderTemplate("env.example.tmpl", project)
	if err != nil {
		return nil, err
	}

	requirements, err := g.renderTemplate("requirements.txt.tmpl", project)
	if err != nil {
		return nil, err
	}

	files := []File{
		{Path: filepath.Join(rootPackage, "__init__.py"), Content: "from . import agent\n"},
		{Path: filepath.Join(rootPackage, "agent.py"), Content: rootAgent},
		{Path: filepath.Join(rootPackage, ".env.example"), Content: envExample},
		{Path: filepath.Join(rootPackage, "sub_agents", "__init__.py"), Content: ""},
		{Path: "requirements.txt", Content: requirements},
	}

	if project.AddReadme {
		readme, err := g.renderTemplate("README.md.tmpl", project)
		if err != nil {
			return nil, err
		}
		files = append(files, File{Path: "README.md", Content: readme})
	}

	for _, agent := range project.Orchestrator.SubAgents {
		agentFiles, err := g.generateAgentPackageFiles(filepath.Join(rootPackage, "sub_agents"), agent)
		if err != nil {
			return nil, err
		}
		files = append(files, agentFiles...)
	}

	return files, nil
}

func (g *Generator) GenerateSingleAgentFiles(agent *model.Agent) ([]File, error) {
	return g.generateAgentPackageFiles(".", agent)
}

func (g *Generator) generateAgentPackageFiles(baseDir string, agent *model.Agent) ([]File, error) {
	content, err := g.renderAgentTemplate(agent)
	if err != nil {
		return nil, err
	}

	packageDir := filepath.Join(baseDir, agent.PackageName())
	return []File{
		{Path: filepath.Join(packageDir, "__init__.py"), Content: "from .agent import agent\n"},
		{Path: filepath.Join(packageDir, "agent.py"), Content: content},
	}, nil
}

func (g *Generator) renderAgentTemplate(agent *model.Agent) (string, error) {
	switch agent.Type {
	case model.AgentTypeLLM:
		return g.renderTemplate("llm_agent.py.tmpl", agent)
	case model.AgentTypeCustom:
		return g.renderTemplate("custom_agent.py.tmpl", agent)
	default:
		return "", fmt.Errorf("unsupported agent type: %s", agent.Type)
	}
}

func (g *Generator) renderTemplate(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := g.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", fmt.Errorf("failed to generate %s: %w", name, err)
	}
	return buf.String(), nil
}

func snakeCase(name string) string {
	return model.NewAgent(name, model.AgentTypeCustom, "", "output", "").PackageName()
}

func className(name string) string {
	return model.NewAgent(name, model.AgentTypeCustom, "", "output", "").ClassName()
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

func getOrchestratorImports(pattern model.OrchestrationPattern) string {
	return getAgentClass(pattern)
}

func coordinatorInstruction(orchestrator *model.Orchestrator) string {
	base := "Coordinate the available sub-agents to solve the user's request. Delegate work to the most appropriate sub-agent and provide a direct final answer to the user."
	if orchestrator.Description == "" {
		return base
	}
	return base + "\n\nProject context: " + strings.TrimSpace(orchestrator.Description)
}
