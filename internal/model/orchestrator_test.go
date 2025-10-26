package model

import (
	"testing"
)

func TestNewOrchestrator(t *testing.T) {
	tests := []struct {
		name        string
		orchName    string
		pattern     OrchestrationPattern
		description string
		model       string
		want        *Orchestrator
	}{
		{
			name:        "create sequential orchestrator",
			orchName:    "Coordinator",
			pattern:     PatternSequential,
			description: "Coordinates tasks",
			model:       "gemini-2.0-flash",
			want: &Orchestrator{
				Name:        "Coordinator",
				Pattern:     PatternSequential,
				Description: "Coordinates tasks",
				Model:       "gemini-2.0-flash",
				SubAgents:   []*Agent{},
			},
		},
		{
			name:        "create parallel orchestrator",
			orchName:    "ParallelCoord",
			pattern:     PatternParallel,
			description: "Runs tasks in parallel",
			model:       "gemini-2.0-flash",
			want: &Orchestrator{
				Name:        "ParallelCoord",
				Pattern:     PatternParallel,
				Description: "Runs tasks in parallel",
				Model:       "gemini-2.0-flash",
				SubAgents:   []*Agent{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOrchestrator(tt.orchName, tt.pattern, tt.description, tt.model)

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Pattern != tt.want.Pattern {
				t.Errorf("Pattern = %v, want %v", got.Pattern, tt.want.Pattern)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.Model != tt.want.Model {
				t.Errorf("Model = %v, want %v", got.Model, tt.want.Model)
			}
			if got.SubAgents == nil {
				t.Error("SubAgents should be initialized as empty slice")
			}
		})
	}
}

func TestOrchestrator_AddSubAgent(t *testing.T) {
	orch := NewOrchestrator("Coordinator", PatternSequential, "Test", "gemini-2.0-flash")
	agent1 := NewAgent("Agent1", AgentTypeLLM, "Do task 1", "result1", "gemini-2.0-flash")
	agent2 := NewAgent("Agent2", AgentTypeLLM, "Do task 2", "result2", "gemini-2.0-flash")

	orch.AddSubAgent(agent1)
	if len(orch.SubAgents) != 1 {
		t.Errorf("SubAgents length = %d, want 1", len(orch.SubAgents))
	}

	orch.AddSubAgent(agent2)
	if len(orch.SubAgents) != 2 {
		t.Errorf("SubAgents length = %d, want 2", len(orch.SubAgents))
	}

	if orch.SubAgents[0].Name != "Agent1" {
		t.Errorf("First agent name = %s, want Agent1", orch.SubAgents[0].Name)
	}
	if orch.SubAgents[1].Name != "Agent2" {
		t.Errorf("Second agent name = %s, want Agent2", orch.SubAgents[1].Name)
	}
}

func TestOrchestrator_Validate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Orchestrator
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid orchestrator with agents",
			setup: func() *Orchestrator {
				orch := NewOrchestrator("Coordinator", PatternSequential, "Test", "gemini-2.0-flash")
				orch.AddSubAgent(NewAgent("Agent1", AgentTypeLLM, "Task", "result", "gemini-2.0-flash"))
				return orch
			},
			wantErr: false,
		},
		{
			name: "empty name returns error",
			setup: func() *Orchestrator {
				return NewOrchestrator("", PatternSequential, "Test", "gemini-2.0-flash")
			},
			wantErr: true,
			errMsg:  "orchestrator name cannot be empty",
		},
		{
			name: "no sub-agents returns error",
			setup: func() *Orchestrator {
				return NewOrchestrator("Coordinator", PatternSequential, "Test", "gemini-2.0-flash")
			},
			wantErr: true,
			errMsg:  "orchestrator must have at least one sub-agent",
		},
		{
			name: "invalid sub-agent returns error",
			setup: func() *Orchestrator {
				orch := NewOrchestrator("Coordinator", PatternSequential, "Test", "gemini-2.0-flash")
				orch.AddSubAgent(NewAgent("", AgentTypeLLM, "Task", "result", "gemini-2.0-flash"))
				return orch
			},
			wantErr: true,
			errMsg:  "sub-agent validation failed: name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orch := tt.setup()
			err := orch.Validate()

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

func TestOrchestrationPattern_String(t *testing.T) {
	tests := []struct {
		pattern OrchestrationPattern
		want    string
	}{
		{PatternSequential, "Sequential"},
		{PatternParallel, "Parallel"},
		{PatternLLMCoordinated, "LLM-Coordinated"},
		{PatternLoop, "Loop"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.pattern.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrchestrationPattern_Description(t *testing.T) {
	tests := []struct {
		pattern OrchestrationPattern
		want    string
	}{
		{PatternSequential, "Sub-agents run one after another"},
		{PatternParallel, "Sub-agents run simultaneously"},
		{PatternLLMCoordinated, "Orchestrator decides which sub-agent to call"},
		{PatternLoop, "Repeat sub-agents until condition met"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern.String(), func(t *testing.T) {
			if got := tt.pattern.Description(); got != tt.want {
				t.Errorf("Description() = %v, want %v", got, tt.want)
			}
		})
	}
}
