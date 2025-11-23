package builtins

import (
	"fmt"
	"io"
	"os"
)

const PwdCommandName = "pwd"

// PwdCommand реализует встроенную команду pwd.
// Выводит текущую рабочую директорию.
type PwdCommand struct{}

// NewPwdCommand создает новый экземпляр команды pwd.
func NewPwdCommand() *PwdCommand {
	return &PwdCommand{}
}

// Name возвращает имя команды pwd.
func (p *PwdCommand) Name() string {
	return PwdCommandName
}

// Execute выполняет команду pwd.
// Выводит абсолютный путь к текущей рабочей директории.
func (p *PwdCommand) Execute(_ []string, _ map[string]string, _ io.Reader, stdout io.Writer, _ io.Writer) int {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, dir)

	return 0
}
