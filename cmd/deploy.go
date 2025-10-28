package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/doji-co/agent-builder/internal/prompt"
	"github.com/doji-co/agent-builder/internal/service"
	"github.com/spf13/cobra"
)

var (
	projectID     string
	region        string
	stagingBucket string
	gcloudService GcloudServiceInterface
	interactive   *prompt.Interactive
)

type GcloudServiceInterface interface {
	IsAvailable() bool
	GetProjectID() (string, error)
	GetRegion() (string, error)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy agent to Vertex AI",
	Long:  "Deploy your multi-agent system to Google Cloud Vertex AI Agent Engine.",
	RunE:  runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&projectID, "project-id", "", "GCP project ID (optional)")
	deployCmd.Flags().StringVar(&region, "region", "", "GCP region (optional)")
	deployCmd.Flags().StringVar(&stagingBucket, "staging-bucket", "", "GCS staging bucket (optional, e.g., gs://my-bucket)")

	gcloudService = service.NewDefaultGcloudService()
	interactive = prompt.NewInteractive()
}

func runDeploy(cmd *cobra.Command, args []string) error {
	deployScript := "./deploy.py"

	if _, err := os.Stat(deployScript); os.IsNotExist(err) {
		return fmt.Errorf("deploy.py not found in current directory. Make sure you're in the project root")
	}

	if projectID == "" && gcloudService.IsAvailable() {
		gcloudProjectID, err := gcloudService.GetProjectID()
		if err == nil && gcloudProjectID != "" {
			useGcloud, err := interactive.PromptUseGcloudValue("project ID", gcloudProjectID)
			if err != nil {
				return fmt.Errorf("failed to prompt for gcloud project ID: %w", err)
			}
			if useGcloud {
				projectID = gcloudProjectID
			}
		}
	}

	if projectID == "" {
		var err error
		projectID, err = interactive.PromptProjectID()
		if err != nil {
			return fmt.Errorf("failed to prompt for project ID: %w", err)
		}
	}

	if region == "" && gcloudService.IsAvailable() {
		gcloudRegion, err := gcloudService.GetRegion()
		if err == nil && gcloudRegion != "" {
			useGcloud, err := interactive.PromptUseGcloudValue("region", gcloudRegion)
			if err != nil {
				return fmt.Errorf("failed to prompt for gcloud region: %w", err)
			}
			if useGcloud {
				region = gcloudRegion
			}
		}
	}

	if region == "" {
		var err error
		region, err = interactive.PromptRegion()
		if err != nil {
			return fmt.Errorf("failed to prompt for region: %w", err)
		}
	}

	if stagingBucket == "" {
		var err error
		stagingBucket, err = interactive.PromptStagingBucket()
		if err != nil {
			return fmt.Errorf("failed to prompt for staging bucket: %w", err)
		}
	}

	fmt.Println("\n🚀 Deploying agent to Vertex AI...")
	fmt.Printf("   Project: %s\n", projectID)
	fmt.Printf("   Region: %s\n", region)
	fmt.Printf("   Bucket: %s\n\n", stagingBucket)

	deployCmd := exec.Command("python3", deployScript,
		"--project-id", projectID,
		"--region", region,
		"--staging-bucket", stagingBucket,
	)

	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr

	if err := deployCmd.Run(); err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	return nil
}
