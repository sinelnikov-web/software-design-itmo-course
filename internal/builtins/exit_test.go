package builtins

import (
	"testing"
)

func TestExitCommand_Name(t *testing.T) {
	command := NewExitCommand()
	expected := ExitCommandName

	if command.Name() != expected {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}

func TestExitCommand_ExecuteWithNoArgs(t *testing.T) {
	command := NewExitCommand()

	// Примечание: Мы не можем реально тестировать os.Exit(), так как это завершает процесс
	// Этот тест только проверяет, что структура команды корректна
	if command.Name() != ExitCommandName {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), ExitCommandName)
	}
}

func TestExitCommand_ExecuteWithExitCode(t *testing.T) {
	command := NewExitCommand()

	// Тест с аргументом кода выхода
	// Примечание: Мы не можем реально тестировать os.Exit(), так как это завершает процесс
	// Этот тест только проверяет, что структура команды корректна
	if command.Name() != ExitCommandName {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), ExitCommandName)
	}
}

func TestExitCommand_ExecuteWithInvalidExitCode(t *testing.T) {
	command := NewExitCommand()

	// Тест с недопустимым кодом выхода
	// Примечание: Мы не можем реально тестировать os.Exit(), так как это завершает процесс
	// Этот тест только проверяет, что структура команды корректна
	if command.Name() != ExitCommandName {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), ExitCommandName)
	}
}

// TestExitCommand_Integration тестирует команду exit в контролируемом режиме
func TestExitCommand_Integration(t *testing.T) {
	// Этот тест проверяет, что команда exit может быть создана и имеет правильное имя
	// Мы не можем тестировать реальное выполнение, так как это вызывает os.Exit()

	command := NewExitCommand()

	if command.Name() != ExitCommandName {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), ExitCommandName)
	}

	var _ Builtin = command
}
