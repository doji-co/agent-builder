package prompt

import (
	"reflect"
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
			name:    "valid with underscores",
			input:   "My_Agent",
			wantErr: false,
		},
		{
			name:    "valid with hyphens",
			input:   "grafana-agent",
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

func TestValidateOutputKey(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid output key",
			input:   "research_notes",
			wantErr: false,
		},
		{
			name:    "empty output key",
			input:   "",
			wantErr: true,
		},
		{
			name:    "output key with hyphen",
			input:   "research-notes",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOutputKey(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResolveModel(t *testing.T) {
	tests := []struct {
		name      string
		selection string
		custom    string
		want      string
	}{
		{
			name:      "uses selected model",
			selection: "gemini-3.1-pro-preview",
			want:      "gemini-3.1-pro-preview",
		},
		{
			name:      "uses custom model value",
			selection: CustomModelOption,
			custom:    "gemini-pro-latest",
			want:      "gemini-pro-latest",
		},
		{
			name:      "falls back to default model for empty custom value",
			selection: CustomModelOption,
			want:      DefaultModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveModel(tt.selection, tt.custom); got != tt.want {
				t.Errorf("ResolveModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAvailableModels(t *testing.T) {
	if DefaultModel != "gemini-2.5-flash" {
		t.Fatalf("DefaultModel = %q, want %q", DefaultModel, "gemini-2.5-flash")
	}

	want := []string{
		DefaultModel,
		"gemini-3-flash-preview",
		"gemini-3.1-pro-preview",
		"gemini-3.1-flash-lite",
		"gemini-2.5-pro",
		"gemini-2.5-flash-lite",
		CustomModelOption,
	}

	if !reflect.DeepEqual(AvailableModels, want) {
		t.Errorf("AvailableModels = %v, want %v", AvailableModels, want)
	}
}
