package builtins

import (
	"bytes"
	"testing"
)

func TestEchoCommand_Execute(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "single argument",
			args:     []string{"hello"},
			expected: "hello\n",
		},
		{
			name:     "multiple arguments",
			args:     []string{"hello", "world"},
			expected: "hello world\n",
		},
		{
			name:     "no arguments",
			args:     []string{},
			expected: "\n",
		},
		{
			name:     "empty string",
			args:     []string{""},
			expected: "\n",
		},
	}

	command := NewEchoCommand()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			exitCode := command.Execute(tt.args, nil, nil, &stdout, &stderr)

			if exitCode != 0 {
				t.Errorf("EchoCommand.Execute() exitCode = %d, expected 0", exitCode)
			}

			if stdout.String() != tt.expected {
				t.Errorf("EchoCommand.Execute() output = %q, expected %q", stdout.String(), tt.expected)
			}
		})
	}
}

func TestEchoCommand_Name(t *testing.T) {
	command := NewEchoCommand()
	expected := "echo"

	if command.Name() != expected {
		t.Errorf("EchoCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}
