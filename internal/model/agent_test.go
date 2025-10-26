package model

import (
	"testing"
)

func TestNewAgent(t *testing.T) {
	tests := []struct {
		name        string
		agentName   string
		agentType   AgentType
		instruction string
		outputKey   string
		model       string
		want        *Agent
	}{
		{
			name:        "create LLM agent with all fields",
			agentName:   "Researcher",
			agentType:   AgentTypeLLM,
			instruction: "Research the topic",
			outputKey:   "research_data",
			model:       "gemini-2.0-flash",
			want: &Agent{
				Name:        "Researcher",
				Type:        AgentTypeLLM,
				Instruction: "Research the topic",
				OutputKey:   "research_data",
				Model:       "gemini-2.0-flash",
			},
		},
		{
			name:        "create custom agent",
			agentName:   "CustomProcessor",
			agentType:   AgentTypeCustom,
			instruction: "",
			outputKey:   "result",
			model:       "",
			want: &Agent{
				Name:        "CustomProcessor",
				Type:        AgentTypeCustom,
				Instruction: "",
				OutputKey:   "result",
				Model:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAgent(tt.agentName, tt.agentType, tt.instruction, tt.outputKey, tt.model)

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Type != tt.want.Type {
				t.Errorf("Type = %v, want %v", got.Type, tt.want.Type)
			}
			if got.Instruction != tt.want.Instruction {
				t.Errorf("Instruction = %v, want %v", got.Instruction, tt.want.Instruction)
			}
			if got.OutputKey != tt.want.OutputKey {
				t.Errorf("OutputKey = %v, want %v", got.OutputKey, tt.want.OutputKey)
			}
			if got.Model != tt.want.Model {
				t.Errorf("Model = %v, want %v", got.Model, tt.want.Model)
			}
		})
	}
}

func TestAgent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		agent   *Agent
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid LLM agent",
			agent: &Agent{
				Name:        "Researcher",
				Type:        AgentTypeLLM,
				Instruction: "Research the topic",
				OutputKey:   "data",
				Model:       "gemini-2.0-flash",
			},
			wantErr: false,
		},
		{
			name: "empty name returns error",
			agent: &Agent{
				Name:        "",
				Type:        AgentTypeLLM,
				Instruction: "Research",
				OutputKey:   "data",
				Model:       "gemini-2.0-flash",
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "LLM agent without instruction returns error",
			agent: &Agent{
				Name:        "Researcher",
				Type:        AgentTypeLLM,
				Instruction: "",
				OutputKey:   "data",
				Model:       "gemini-2.0-flash",
			},
			wantErr: true,
			errMsg:  "instruction is required for LLM agents",
		},
		{
			name: "custom agent can have empty instruction",
			agent: &Agent{
				Name:      "CustomAgent",
				Type:      AgentTypeCustom,
				OutputKey: "result",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.agent.Validate()

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
