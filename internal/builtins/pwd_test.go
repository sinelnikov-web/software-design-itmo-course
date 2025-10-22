package builtins

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

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

func TestPwdCommand_Name(t *testing.T) {
	command := NewPwdCommand()
	expected := "pwd"

	if command.Name() != expected {
		t.Errorf("PwdCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}

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
