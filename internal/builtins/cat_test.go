package builtins

import (
	"bytes"
	"strings"
	"testing"
)

// TestCatCommand_Execute тестирует выполнение команды cat с различными сценариями.
// Проверяет чтение из stdin при отсутствии аргументов и обработку ошибок
// при попытке чтения несуществующих файлов.
func TestCatCommand_Execute(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		wantErr  bool
	}{
		{
			name:     "no arguments - read from stdin",
			args:     []string{},
			expected: "",
			wantErr:  false,
		},
		{
			name:     "single file",
			args:     []string{"test.txt"},
			expected: "",
			wantErr:  true, // File doesn't exist
		},
		{
			name:     "multiple files",
			args:     []string{"file1.txt", "file2.txt"},
			expected: "",
			wantErr:  true, // Files don't exist
		},
	}

	command := NewCatCommand()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			var stdin bytes.Buffer

			if tt.name == "no arguments - read from stdin" {
				stdin.WriteString("test input")
			}

			exitCode := command.Execute(tt.args, nil, &stdin, &stdout, &stderr)

			if tt.wantErr {
				if exitCode == 0 {
					t.Errorf("CatCommand.Execute() expected error, got exit code 0")
				}
			} else {
				if exitCode != 0 {
					t.Errorf("CatCommand.Execute() exitCode = %d, expected 0", exitCode)
				}
			}
		})
	}
}

// TestCatCommand_ExecuteWithStdin тестирует выполнение команды cat при чтении из stdin.
// Проверяет, что команда корректно читает и выводит данные из стандартного ввода.
func TestCatCommand_ExecuteWithStdin(t *testing.T) {
	command := NewCatCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("hello world\n")

	exitCode := command.Execute([]string{}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("CatCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "hello world") {
		t.Errorf("CatCommand.Execute() output = %q, expected to contain 'hello world'", output)
	}
}

// TestCatCommand_Name тестирует получение имени команды cat.
// Проверяет, что команда возвращает корректное имя "cat".
func TestCatCommand_Name(t *testing.T) {
	command := NewCatCommand()
	expected := "cat"

	if command.Name() != expected {
		t.Errorf("CatCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}
