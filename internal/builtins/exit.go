package builtins

import (
	"io"
	"os"
	"strconv"
)

const ExitCommandName = "exit"

// ExitCommand реализует встроенную команду exit.
// Завершает работу shell'а с указанным кодом возврата.
type ExitCommand struct{}

// NewExitCommand создает новый экземпляр команды exit.
func NewExitCommand() *ExitCommand {
	return &ExitCommand{}
}

// Name возвращает имя команды exit.
func (e *ExitCommand) Name() string {
	return ExitCommandName
}

// Execute выполняет команду exit.
// Если передан аргумент, использует его как код возврата.
// Иначе завершает работу с кодом 0.
func (e *ExitCommand) Execute(args []string, _ map[string]string, _ io.Reader, _ io.Writer, _ io.Writer) int {
	exitCode := 0

	if len(args) > 0 {
		if code, err := strconv.Atoi(args[0]); err == nil {
			exitCode = code
		}
	}

	os.Exit(exitCode)

	return exitCode
}
