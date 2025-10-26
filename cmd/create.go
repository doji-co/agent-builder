package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/doji-co/agent-builder/internal/generator"
	"github.com/doji-co/agent-builder/internal/model"
	"github.com/doji-co/agent-builder/internal/prompt"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new multi-agent project",
	Long:  "Launch an interactive session to create a new ADK multi-agent project.",
	RunE:  runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	interactive := prompt.NewInteractive()

	fmt.Println("🤖 Welcome to Agent Builder!")
	fmt.Println("Let's create your multi-agent system.")

	projectName, err := interactive.PromptProjectName()
	if err != nil {
		return fmt.Errorf("failed to get project name: %w", err)
	}

	pattern, err := interactive.PromptOrchestrationPattern()
	if err != nil {
		return fmt.Errorf("failed to get orchestration pattern: %w", err)
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 ORCHESTRATOR CONFIGURATION")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	orchName, err := interactive.PromptOrchestratorName()
	if err != nil {
		return fmt.Errorf("failed to get orchestrator name: %w", err)
	}

	orchDescription, err := interactive.PromptOrchestratorDescription()
	if err != nil {
		return fmt.Errorf("failed to get orchestrator description: %w", err)
	}

	orchModel, err := interactive.PromptModel(prompt.DefaultModel)
	if err != nil {
		return fmt.Errorf("failed to get orchestrator model: %w", err)
	}

	orchestrator := model.NewOrchestrator(orchName, pattern, orchDescription, orchModel)

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🤖 SUB-AGENTS CONFIGURATION")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	agentNumber := 1
	for {
		agentName, err := interactive.PromptAgentName(agentNumber)
		if err != nil {
			return fmt.Errorf("failed to get agent name: %w", err)
		}

		agentType, err := interactive.PromptAgentType()
		if err != nil {
			return fmt.Errorf("failed to get agent type: %w", err)
		}

		var instruction string
		if agentType == model.AgentTypeLLM {
			instruction, err = interactive.PromptAgentInstruction(agentName)
			if err != nil {
				return fmt.Errorf("failed to get agent instruction: %w", err)
			}
		}

		outputKey, err := interactive.PromptOutputKey()
		if err != nil {
			return fmt.Errorf("failed to get output key: %w", err)
		}

		agentModel, err := interactive.PromptModel(prompt.DefaultModel)
		if err != nil {
			return fmt.Errorf("failed to get agent model: %w", err)
		}

		agent := model.NewAgent(agentName, agentType, instruction, outputKey, agentModel)
		orchestrator.AddSubAgent(agent)

		fmt.Printf("\n✓ Sub-agent \"%s\" added to %s\n\n", agentName, orchName)

		addMore, err := interactive.PromptAddAnotherAgent()
		if err != nil {
			return fmt.Errorf("failed to prompt for another agent: %w", err)
		}

		if !addMore {
			break
		}

		agentNumber++
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📦 PROJECT SETUP")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	project := model.NewProject(projectName, orchestrator)

	outputDir, err := interactive.PromptOutputDirectory(project.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to get output directory: %w", err)
	}
	project.OutputDir = outputDir

	addExample, err := interactive.PromptAddExample()
	if err != nil {
		return fmt.Errorf("failed to prompt for example: %w", err)
	}
	project.AddExample = addExample

	addDocker, err := interactive.PromptAddDocker()
	if err != nil {
		return fmt.Errorf("failed to prompt for Docker: %w", err)
	}
	project.AddDocker = addDocker

	if err := project.Validate(); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	fmt.Println("\n✨ Generating project structure...")

	if err := generateProject(project); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Println("\n📁 System Architecture:")
	fmt.Printf("   %s (%s)\n", orchestrator.Name, pattern.String())
	for _, agent := range orchestrator.SubAgents {
		fmt.Printf("   ├── %s (%s)\n", agent.Name, agent.Type)
	}

	fmt.Printf("\n✓ Created %s/\n", project.OutputDir)
	fmt.Println("  ├── agents.py          # Agent definitions")
	if project.AddExample {
		fmt.Println("  ├── main.py            # Example usage")
	}
	fmt.Println("  ├── requirements.txt   # Dependencies (google-adk)")
	if project.AddReadme {
		fmt.Println("  └── README.md          # Documentation")
	}

	fmt.Println("\n🚀 Next steps:")
	fmt.Printf("  cd %s\n", project.OutputDir)
	fmt.Println("  pip install -r requirements.txt")
	if project.AddExample {
		fmt.Println("  python main.py \"Your prompt here\"")
	}

	return nil
}

func generateProject(project *model.Project) error {
	if err := os.MkdirAll(project.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	gen := generator.NewGenerator()

	agentsPy, err := gen.GenerateAgentsPy(project)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(project.OutputDir, "agents.py"), []byte(agentsPy), 0644); err != nil {
		return fmt.Errorf("failed to write agents.py: %w", err)
	}

	if project.AddExample {
		mainPy, err := gen.GenerateMainPy(project)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(project.OutputDir, "main.py"), []byte(mainPy), 0644); err != nil {
			return fmt.Errorf("failed to write main.py: %w", err)
		}
	}

	requirementsTxt, err := gen.GenerateRequirementsTxt()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(project.OutputDir, "requirements.txt"), []byte(requirementsTxt), 0644); err != nil {
		return fmt.Errorf("failed to write requirements.txt: %w", err)
	}

	if project.AddReadme {
		readme, err := gen.GenerateReadme(project)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(project.OutputDir, "README.md"), []byte(readme), 0644); err != nil {
			return fmt.Errorf("failed to write README.md: %w", err)
		}
	}

	return nil
}
