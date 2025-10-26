package model

import (
	"testing"
)

func TestNewProject(t *testing.T) {
	orch := NewOrchestrator("Coordinator", PatternSequential, "Test coordinator", "gemini-2.0-flash")
	orch.AddSubAgent(NewAgent("Agent1", AgentTypeLLM, "Task 1", "result1", "gemini-2.0-flash"))

	project := NewProject("my-project", orch)

	if project.Name != "my-project" {
		t.Errorf("Name = %v, want my-project", project.Name)
	}
	if project.Orchestrator != orch {
		t.Errorf("Orchestrator not set correctly")
	}
	if project.OutputDir == "" {
		t.Error("OutputDir should have default value")
	}
}

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Project
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid project",
			setup: func() *Project {
				orch := NewOrchestrator("Coordinator", PatternSequential, "Test", "gemini-2.0-flash")
				orch.AddSubAgent(NewAgent("Agent1", AgentTypeLLM, "Task", "result", "gemini-2.0-flash"))
				return NewProject("my-project", orch)
			},
			wantErr: false,
		},
		{
			name: "empty name returns error",
			setup: func() *Project {
				orch := NewOrchestrator("Coordinator", PatternSequential, "Test", "gemini-2.0-flash")
				orch.AddSubAgent(NewAgent("Agent1", AgentTypeLLM, "Task", "result", "gemini-2.0-flash"))
				return NewProject("", orch)
			},
			wantErr: true,
			errMsg:  "project name cannot be empty",
		},
		{
			name: "nil orchestrator returns error",
			setup: func() *Project {
				return NewProject("my-project", nil)
			},
			wantErr: true,
			errMsg:  "orchestrator cannot be nil",
		},
		{
			name: "invalid orchestrator returns error",
			setup: func() *Project {
				orch := NewOrchestrator("", PatternSequential, "Test", "gemini-2.0-flash")
				return NewProject("my-project", orch)
			},
			wantErr: true,
			errMsg:  "orchestrator validation failed: orchestrator name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := tt.setup()
			err := project.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
