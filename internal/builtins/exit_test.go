package builtins

import (
	"testing"
)

// TestExitCommand_Name тестирует получение имени команды exit.
// Проверяет, что команда возвращает корректное имя "exit".
func TestExitCommand_Name(t *testing.T) {
	command := NewExitCommand()
	expected := ExitCommandName

	if command.Name() != expected {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}

// TestExitCommand_ExecuteWithNoArgs тестирует команду exit без аргументов.
// Примечание: Мы не можем реально тестировать os.Exit(), так как это завершает процесс.
// Этот тест только проверяет, что структура команды корректна.
func TestExitCommand_ExecuteWithNoArgs(t *testing.T) {
	command := NewExitCommand()
	if command.Name() != ExitCommandName {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), ExitCommandName)
	}
}

// TestExitCommand_ExecuteWithExitCode тестирует команду exit с аргументом кода выхода.
// Примечание: Мы не можем реально тестировать os.Exit(), так как это завершает процесс.
// Этот тест только проверяет, что структура команды корректна.
func TestExitCommand_ExecuteWithExitCode(t *testing.T) {
	command := NewExitCommand()
	if command.Name() != ExitCommandName {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), ExitCommandName)
	}
}

// TestExitCommand_ExecuteWithInvalidExitCode тестирует команду exit с недопустимым кодом выхода.
// Примечание: Мы не можем реально тестировать os.Exit(), так как это завершает процесс.
// Этот тест только проверяет, что структура команды корректна.
//
// Ожидаемое поведение при невалидном аргументе:
//   - Команда выводит ошибку в stderr: "exit: <arg>: numeric argument required"
//   - Завершает процесс с кодом 2
//
// Ожидаемое поведение при коде вне диапазона 0-255:
//   - Код приводится к диапазону 0-255 (по модулю 256)
//   - Например: exit 256 -> код 0, exit -1 -> код 255, exit 300 -> код 44
func TestExitCommand_ExecuteWithInvalidExitCode(t *testing.T) {
	command := NewExitCommand()
	if command.Name() != ExitCommandName {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), ExitCommandName)
	}
}

// TestExitCommand_Integration тестирует команду exit в контролируемом режиме.
// Проверяет, что команда exit может быть создана и имеет правильное имя.
// Мы не можем тестировать реальное выполнение, так как это вызывает os.Exit().
func TestExitCommand_Integration(t *testing.T) {

	command := NewExitCommand()

	if command.Name() != ExitCommandName {
		t.Errorf("ExitCommand.Name() = %s, expected %s", command.Name(), ExitCommandName)
	}

	var _ Builtin = command
}
