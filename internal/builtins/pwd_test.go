package builtins

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestPwdCommand_Execute тестирует выполнение команды pwd без аргументов.
// Проверяет, что команда возвращает текущую рабочую директорию.
func TestPwdCommand_Execute(t *testing.T) {
	command := NewPwdCommand()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := command.Execute([]string{}, nil, nil, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("PwdCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := strings.TrimSpace(stdout.String())
	expected, _ := os.Getwd()

	if output != expected {
		t.Errorf("PwdCommand.Execute() output = %q, expected %q", output, expected)
	}
}

// TestPwdCommand_Name тестирует получение имени команды pwd.
// Проверяет, что команда возвращает корректное имя "pwd".
func TestPwdCommand_Name(t *testing.T) {
	command := NewPwdCommand()
	expected := "pwd"

	if command.Name() != expected {
		t.Errorf("PwdCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}

// TestPwdCommand_ExecuteWithArgs тестирует выполнение команды pwd с аргументами.
// Проверяет, что команда игнорирует аргументы и все равно возвращает текущую директорию.
func TestPwdCommand_ExecuteWithArgs(t *testing.T) {
	command := NewPwdCommand()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := command.Execute([]string{"ignored", "args"}, nil, nil, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("PwdCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := strings.TrimSpace(stdout.String())
	expected, _ := os.Getwd()

	if output != expected {
		t.Errorf("PwdCommand.Execute() output = %q, expected %q", output, expected)
	}
}
