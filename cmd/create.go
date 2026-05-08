package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/doji-co/agent-builder/internal/generator"
	"github.com/doji-co/agent-builder/internal/model"
	"github.com/doji-co/agent-builder/internal/prompt"
	"github.com/spf13/cobra"
)

type createPrompter interface {
	PromptProjectType() (string, error)
	PromptProjectName() (string, error)
	PromptOrchestrationPattern() (model.OrchestrationPattern, error)
	PromptOrchestratorName() (string, error)
	PromptOrchestratorDescription() (string, error)
	PromptModel(defaultModel string) (string, error)
	PromptAgentName(agentNumber int) (string, error)
	PromptAgentType() (model.AgentType, error)
	PromptAgentInstruction(agentName string) (string, error)
	PromptOutputKey() (string, error)
	PromptAddAnotherAgent() (bool, error)
	PromptOutputDirectory(defaultDir string) (string, error)
}

type scaffoldGenerator interface {
	GenerateProjectFiles(project *model.Project) ([]generator.File, error)
	GenerateSingleAgentFiles(agent *model.Agent) ([]generator.File, error)
}

type fileWriter interface {
	WriteFiles(rootDir string, files []generator.File) error
}

type createDependencies struct {
	prompterFactory  func() createPrompter
	generatorFactory func() scaffoldGenerator
	fileWriter       fileWriter
}

func (d createDependencies) withDefaults() createDependencies {
	if d.prompterFactory == nil {
		d.prompterFactory = func() createPrompter { return prompt.NewInteractive() }
	}
	if d.generatorFactory == nil {
		d.generatorFactory = func() scaffoldGenerator { return generator.NewGenerator() }
	}
	if d.fileWriter == nil {
		d.fileWriter = osFileWriter{}
	}
	return d
}

type osFileWriter struct{}

func (osFileWriter) WriteFiles(rootDir string, files []generator.File) error {
	for _, file := range files {
		targetPath := filepath.Join(rootDir, file.Path)
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(targetPath, []byte(file.Content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func newCreateCommand(deps createDependencies) *cobra.Command {
	deps = deps.withDefaults()

	return &cobra.Command{
		Use:   "create",
		Short: "Create a new multi-agent project",
		Long:  "Launch an interactive session to create a new ADK multi-agent project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(cmd, deps)
		},
	}
}

func runCreate(cmd *cobra.Command, deps createDependencies) error {
	prompter := deps.prompterFactory()
	gen := deps.generatorFactory()
	out := cmd.OutOrStdout()

	fmt.Fprintln(out, "Welcome to Agent Builder")

	projectType, err := prompter.PromptProjectType()
	if err != nil {
		return fmt.Errorf("failed to get project type: %w", err)
	}

	if projectType == "full" {
		return runCreateFullProject(out, prompter, gen, deps.fileWriter)
	}
	return runCreateSingleAgent(out, prompter, gen, deps.fileWriter)
}

func runCreateFullProject(out io.Writer, prompter createPrompter, gen scaffoldGenerator, writer fileWriter) error {
	projectName, err := prompter.PromptProjectName()
	if err != nil {
		return fmt.Errorf("failed to get project name: %w", err)
	}

	pattern, err := prompter.PromptOrchestrationPattern()
	if err != nil {
		return fmt.Errorf("failed to get orchestration pattern: %w", err)
	}

	orchestratorName, err := prompter.PromptOrchestratorName()
	if err != nil {
		return fmt.Errorf("failed to get orchestrator name: %w", err)
	}

	orchestratorDescription, err := prompter.PromptOrchestratorDescription()
	if err != nil {
		return fmt.Errorf("failed to get orchestrator description: %w", err)
	}

	orchestratorModel, err := prompter.PromptModel(prompt.DefaultModel)
	if err != nil {
		return fmt.Errorf("failed to get orchestrator model: %w", err)
	}

	orchestrator := model.NewOrchestrator(orchestratorName, pattern, orchestratorDescription, orchestratorModel)

	agentNumber := 1
	for {
		agent, err := promptForAgent(prompter, agentNumber)
		if err != nil {
			return err
		}
		orchestrator.AddSubAgent(agent)

		addAnother, err := prompter.PromptAddAnotherAgent()
		if err != nil {
			return fmt.Errorf("failed to prompt for another agent: %w", err)
		}
		if !addAnother {
			break
		}
		agentNumber++
	}

	project := model.NewProject(projectName, orchestrator)
	outputDir, err := prompter.PromptOutputDirectory(project.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to get output directory: %w", err)
	}
	project.OutputDir = outputDir

	if err := project.Validate(); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	files, err := gen.GenerateProjectFiles(project)
	if err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
	}

	if err := writer.WriteFiles(project.OutputDir, files); err != nil {
		return fmt.Errorf("failed to write project files: %w", err)
	}

	fmt.Fprintf(out, "Created %s\n", project.OutputDir)
	fmt.Fprintf(out, "Run from %s:\n", project.OutputDir)
	fmt.Fprintln(out, "python3 -m venv .venv")
	fmt.Fprintln(out, "source .venv/bin/activate")
	fmt.Fprintln(out, "pip install -r requirements.txt")
	fmt.Fprintf(out, "cp %s/.env.example %s/.env\n", project.PackageName(), project.PackageName())
	fmt.Fprintf(out, "adk run %s\n", project.PackageName())
	fmt.Fprintln(out, "adk web --port 8000")

	return nil
}

func runCreateSingleAgent(out io.Writer, prompter createPrompter, gen scaffoldGenerator, writer fileWriter) error {
	agent, err := promptForAgent(prompter, 1)
	if err != nil {
		return err
	}

	if err := agent.Validate(); err != nil {
		return fmt.Errorf("agent validation failed: %w", err)
	}

	files, err := gen.GenerateSingleAgentFiles(agent)
	if err != nil {
		return fmt.Errorf("failed to generate agent files: %w", err)
	}

	if err := writer.WriteFiles(".", files); err != nil {
		return fmt.Errorf("failed to write agent files: %w", err)
	}

	fmt.Fprintf(out, "Created %s/\n", agent.PackageName())
	fmt.Fprintf(out, "Import this package from your root ADK agent and add it to sub_agents.\n")
	return nil
}

func promptForAgent(prompter createPrompter, agentNumber int) (*model.Agent, error) {
	agentName, err := prompter.PromptAgentName(agentNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent name: %w", err)
	}

	agentType, err := prompter.PromptAgentType()
	if err != nil {
		return nil, fmt.Errorf("failed to get agent type: %w", err)
	}

	var instruction string
	if agentType == model.AgentTypeLLM {
		instruction, err = prompter.PromptAgentInstruction(agentName)
		if err != nil {
			return nil, fmt.Errorf("failed to get agent instruction: %w", err)
		}
	}

	outputKey, err := prompter.PromptOutputKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get output key: %w", err)
	}

	var modelName string
	if agentType == model.AgentTypeLLM {
		modelName, err = prompter.PromptModel(prompt.DefaultModel)
		if err != nil {
			return nil, fmt.Errorf("failed to get agent model: %w", err)
		}
	}

	agent := model.NewAgent(agentName, agentType, instruction, outputKey, modelName)
	if err := agent.Validate(); err != nil {
		return nil, fmt.Errorf("agent validation failed: %w", err)
	}

	return agent, nil
}
