package builtins

import (
	"bytes"
	"testing"
)

// TestEchoCommand_Execute тестирует выполнение команды echo с различными аргументами.
// Проверяет корректность вывода для одного аргумента, нескольких аргументов,
// отсутствия аргументов и пустой строки.
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

// TestEchoCommand_Name тестирует получение имени команды echo.
// Проверяет, что команда возвращает корректное имя "echo".
func TestEchoCommand_Name(t *testing.T) {
	command := NewEchoCommand()
	expected := "echo"

	if command.Name() != expected {
		t.Errorf("EchoCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}
