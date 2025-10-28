package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestDeployCommand_FlagsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "all flags are now optional",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "only project-id provided",
			args:    []string{"--project-id", "my-project"},
			wantErr: false,
		},
		{
			name:    "only staging-bucket provided",
			args:    []string{"--staging-bucket", "gs://bucket"},
			wantErr: false,
		},
		{
			name:    "all flags provided",
			args:    []string{"--project-id", "my-project", "--region", "us-west1", "--staging-bucket", "gs://bucket"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use: "deploy",
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			}

			cmd.Flags().StringVar(&projectID, "project-id", "", "GCP project ID (optional)")
			cmd.Flags().StringVar(&region, "region", "", "GCP region (optional)")
			cmd.Flags().StringVar(&stagingBucket, "staging-bucket", "", "GCS staging bucket (optional)")

			cmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none")
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
		})
	}
}

func TestDeployCommand_DeployScriptNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	var testProjectID, testRegion, testStagingBucket string

	cmd := &cobra.Command{
		Use:  "deploy",
		RunE: runDeploy,
	}

	cmd.Flags().StringVar(&testProjectID, "project-id", "", "GCP project ID (required)")
	cmd.Flags().StringVar(&testRegion, "region", "us-central1", "GCP region")
	cmd.Flags().StringVar(&testStagingBucket, "staging-bucket", "", "GCS staging bucket (required)")

	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("staging-bucket")

	cmd.SetArgs([]string{
		"--project-id", "test-project",
		"--staging-bucket", "gs://test-bucket",
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err = cmd.Execute()

	if err == nil {
		t.Error("expected error when deploy.py not found, got none")
		return
	}

	if !strings.Contains(err.Error(), "deploy.py not found") {
		t.Errorf("error = %v, want error containing 'deploy.py not found'", err)
	}
}

func TestDeployCommand_DeployScriptExists(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(originalDir)

	deployScript := filepath.Join(tmpDir, "deploy.py")
	content := `#!/usr/bin/env python3
import sys
print("Mock deploy script")
sys.exit(0)
`
	if err := os.WriteFile(deployScript, []byte(content), 0755); err != nil {
		t.Fatalf("failed to create deploy.py: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	var testProjectID, testRegion, testStagingBucket string

	cmd := &cobra.Command{
		Use:  "deploy",
		RunE: runDeploy,
	}

	cmd.Flags().StringVar(&testProjectID, "project-id", "", "GCP project ID (required)")
	cmd.Flags().StringVar(&testRegion, "region", "us-central1", "GCP region")
	cmd.Flags().StringVar(&testStagingBucket, "staging-bucket", "", "GCS staging bucket (required)")

	cmd.MarkFlagRequired("project-id")
	cmd.MarkFlagRequired("staging-bucket")

	cmd.SetArgs([]string{
		"--project-id", "test-project",
		"--staging-bucket", "gs://test-bucket",
	})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err = cmd.Execute()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

type mockGcloudService struct {
	isAvailable  bool
	projectID    string
	projectIDErr error
	region       string
	regionErr    error
}

func (m *mockGcloudService) IsAvailable() bool {
	return m.isAvailable
}

func (m *mockGcloudService) GetProjectID() (string, error) {
	return m.projectID, m.projectIDErr
}

func (m *mockGcloudService) GetRegion() (string, error) {
	return m.region, m.regionErr
}
