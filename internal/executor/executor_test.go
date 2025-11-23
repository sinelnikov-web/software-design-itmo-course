package executor

import (
	"testing"

	"gocli/internal/parser"
)

func TestExecutor_IsBuiltin(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"builtin command", "echo", true},
		{"builtin command", "cat", true},
		{"non-builtin command", "ls", false},
		{"empty command", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.IsBuiltin(tt.command)
			if result != tt.expected {
				t.Errorf("Executor.IsBuiltin(%s) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestExecutor_ListBuiltins(t *testing.T) {
	executor := NewExecutor()
	commands := executor.ListBuiltins()

	expectedCount := 5 // cat, echo, wc, pwd, exit
	if len(commands) != expectedCount {
		t.Errorf("Executor.ListBuiltins() returned %d commands, expected %d", len(commands), expectedCount)
	}
}

func TestExecutor_ExecuteCommand(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name    string
		command *parser.Command
		wantErr bool
	}{
		{
			name: "echo command",
			command: &parser.Command{
				Name: "echo",
				Args: []*parser.Argument{{Value: "hello", Quoted: false}},
			},
			wantErr: false,
		},
		{
			name: "pwd command",
			command: &parser.Command{
				Name: "pwd",
				Args: []*parser.Argument{},
			},
			wantErr: false,
		},
		{
			name: "nonexistent command",
			command: &parser.Command{
				Name: "nonexistent",
				Args: []*parser.Argument{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.Execute(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Executor.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecutor_ExecutePipeline(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name     string
		pipeline *parser.Pipeline
		wantErr  bool
	}{
		{
			name: "single command pipeline",
			pipeline: &parser.Pipeline{
				Commands: []*parser.Command{
					{Name: "echo", Args: []*parser.Argument{{Value: "hello", Quoted: false}}},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple commands pipeline",
			pipeline: &parser.Pipeline{
				Commands: []*parser.Command{
					{Name: "echo", Args: []*parser.Argument{{Value: "hello", Quoted: false}}},
					{Name: "wc", Args: []*parser.Argument{}},
				},
			},
			wantErr: true, // Pipeline execution not implemented in Phase 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.Execute(tt.pipeline)

			if (err != nil) != tt.wantErr {
				t.Errorf("Executor.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
