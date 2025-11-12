package builtins

import (
	"fmt"
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
//
// Поведение:
//   - Без аргументов: завершает с кодом 0 (успех)
//   - С числовым аргументом: завершает с указанным кодом (0-255)
//   - С невалидным аргументом: выводит ошибку в stderr и завершает с кодом 2
//   - Коды выхода вне диапазона 0-255: приводятся к диапазону 0-255 (по модулю 256)
func (e *ExitCommand) Execute(args []string, _ map[string]string, _ io.Reader, _ io.Writer, stderr io.Writer) int {
	exitCode := 0

	if len(args) > 0 {
		code, err := strconv.Atoi(args[0])
		if err != nil {
			// Невалидный аргумент: выводим ошибку и завершаем с кодом 2
			fmt.Fprintf(stderr, "exit: %s: numeric argument required\n", args[0])
			os.Exit(2)
			return 2
		}

		// Приводим код к диапазону 0-255 (стандартный диапазон кодов выхода)
		// Отрицательные числа и числа > 255 приводятся к диапазону
		exitCode = code % 256
		if exitCode < 0 {
			exitCode += 256
		}
	}

	os.Exit(exitCode)

	return exitCode
}
