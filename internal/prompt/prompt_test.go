package prompt

import (
	"testing"

	"github.com/doji-co/agent-builder/internal/model"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid project name",
			input:   "my-project",
			wantErr: false,
		},
		{
			name:    "valid with underscores",
			input:   "my_project",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			input:   "project123",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
		},
		{
			name:    "name with spaces",
			input:   "my project",
			wantErr: true,
		},
		{
			name:    "name with special chars",
			input:   "my-project!",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAgentName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid agent name",
			input:   "Researcher",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			input:   "Agent1",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
		},
		{
			name:    "name with spaces",
			input:   "My Agent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAgentName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAgentName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetOrchestrationPatterns(t *testing.T) {
	patterns := GetOrchestrationPatterns()

	if len(patterns) != 4 {
		t.Errorf("Expected 4 patterns, got %d", len(patterns))
	}

	expectedPatterns := []model.OrchestrationPattern{
		model.PatternSequential,
		model.PatternParallel,
		model.PatternLLMCoordinated,
		model.PatternLoop,
	}

	for i, pattern := range expectedPatterns {
		if patterns[i] != pattern {
			t.Errorf("Pattern %d = %v, want %v", i, patterns[i], pattern)
		}
	}
}

func TestGetAgentTypes(t *testing.T) {
	types := GetAgentTypes()

	if len(types) != 2 {
		t.Errorf("Expected 2 agent types, got %d", len(types))
	}

	expectedTypes := []model.AgentType{
		model.AgentTypeLLM,
		model.AgentTypeCustom,
	}

	for i, agentType := range expectedTypes {
		if types[i] != agentType {
			t.Errorf("AgentType %d = %v, want %v", i, types[i], agentType)
		}
	}
}
