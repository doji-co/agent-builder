package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	projectID     string
	region        string
	agentName     string
	stagingBucket string
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy agent to Vertex AI",
	Long:  "Deploy your multi-agent system to Google Cloud Vertex AI Agent Engine.",
	RunE:  runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&projectID, "project-id", "", "GCP project ID (required)")
	deployCmd.Flags().StringVar(&region, "region", "us-central1", "GCP region")
	deployCmd.Flags().StringVar(&agentName, "agent-name", "", "Agent display name (required)")
	deployCmd.Flags().StringVar(&stagingBucket, "staging-bucket", "", "GCS staging bucket (required, e.g., gs://my-bucket)")

	deployCmd.MarkFlagRequired("project-id")
	deployCmd.MarkFlagRequired("agent-name")
	deployCmd.MarkFlagRequired("staging-bucket")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	deployScript := "./deploy.py"

	if _, err := os.Stat(deployScript); os.IsNotExist(err) {
		return fmt.Errorf("deploy.py not found in current directory. Make sure you're in the project root")
	}

	fmt.Println("🚀 Deploying agent to Vertex AI...")
	fmt.Printf("   Project: %s\n", projectID)
	fmt.Printf("   Region: %s\n", region)
	fmt.Printf("   Agent: %s\n", agentName)
	fmt.Printf("   Bucket: %s\n\n", stagingBucket)

	deployCmd := exec.Command("python3", deployScript,
		"--project-id", projectID,
		"--region", region,
		"--agent-name", agentName,
		"--staging-bucket", stagingBucket,
	)

	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr

	if err := deployCmd.Run(); err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	return nil
}
