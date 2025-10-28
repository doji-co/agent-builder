package service

import (
	"fmt"
	"testing"
)

type mockCommandExecutor struct {
	output string
	err    error
}

func (m *mockCommandExecutor) Execute(name string, args ...string) (string, error) {
	return m.output, m.err
}

func TestGcloudService_IsAvailable(t *testing.T) {
	tests := []struct {
		name       string
		mockOutput string
		mockErr    error
		want       bool
	}{
		{
			name:       "gcloud is available",
			mockOutput: "Google Cloud SDK 450.0.0\n",
			mockErr:    nil,
			want:       true,
		},
		{
			name:       "gcloud not found",
			mockOutput: "",
			mockErr:    fmt.Errorf("executable file not found"),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCommandExecutor{
				output: tt.mockOutput,
				err:    tt.mockErr,
			}

			service := NewGcloudService(mock)
			got := service.IsAvailable()

			if got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGcloudService_GetProjectID(t *testing.T) {
	tests := []struct {
		name       string
		mockOutput string
		mockErr    error
		want       string
		wantErr    bool
	}{
		{
			name:       "successful retrieval",
			mockOutput: "my-gcp-project\n",
			mockErr:    nil,
			want:       "my-gcp-project",
			wantErr:    false,
		},
		{
			name:       "gcloud not configured",
			mockOutput: "",
			mockErr:    fmt.Errorf("gcloud command failed"),
			want:       "",
			wantErr:    true,
		},
		{
			name:       "empty output",
			mockOutput: "\n",
			mockErr:    nil,
			want:       "",
			wantErr:    false,
		},
		{
			name:       "output with whitespace",
			mockOutput: "  my-project  \n",
			mockErr:    nil,
			want:       "my-project",
			wantErr:    false,
		},
		{
			name:       "unset value returns empty string",
			mockOutput: "(unset)",
			mockErr:    nil,
			want:       "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCommandExecutor{
				output: tt.mockOutput,
				err:    tt.mockErr,
			}

			service := NewGcloudService(mock)
			got, err := service.GetProjectID()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("GetProjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGcloudService_GetRegion(t *testing.T) {
	tests := []struct {
		name       string
		mockOutput string
		mockErr    error
		want       string
		wantErr    bool
	}{
		{
			name:       "successful retrieval",
			mockOutput: "us-central1\n",
			mockErr:    nil,
			want:       "us-central1",
			wantErr:    false,
		},
		{
			name:       "gcloud not configured",
			mockOutput: "",
			mockErr:    fmt.Errorf("gcloud command failed"),
			want:       "",
			wantErr:    true,
		},
		{
			name:       "region not set returns empty",
			mockOutput: "\n",
			mockErr:    nil,
			want:       "",
			wantErr:    false,
		},
		{
			name:       "output with whitespace",
			mockOutput: "  europe-west1  \n",
			mockErr:    nil,
			want:       "europe-west1",
			wantErr:    false,
		},
		{
			name:       "unset value returns empty string",
			mockOutput: "(unset)",
			mockErr:    nil,
			want:       "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCommandExecutor{
				output: tt.mockOutput,
				err:    tt.mockErr,
			}

			service := NewGcloudService(mock)
			got, err := service.GetRegion()

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("GetRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}
