package builtins

import (
	"fmt"
	"io"
	"strings"
)

const EchoCommandName = "echo"

// EchoCommand реализует встроенную команду echo.
// Выводит переданные аргументы в стандартный поток вывода.
type EchoCommand struct{}

// NewEchoCommand создает новый экземпляр команды echo.
func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

// Name возвращает имя команды echo.
func (e *EchoCommand) Name() string {
	return EchoCommandName
}

// Execute выполняет команду echo.
// Объединяет все аргументы в одну строку и выводит результат.
// Всегда возвращает код успеха (0).
func (e *EchoCommand) Execute(args []string, _ map[string]string, _ io.Reader, stdout io.Writer, _ io.Writer) int {
	output := strings.Join(args, " ")
	fmt.Fprintln(stdout, output)
	return 0
}
