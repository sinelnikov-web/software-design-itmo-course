package builtins

import (
	"fmt"
	"io"
	"os"
)

const CatCommandName = "cat"

// CatCommand реализует встроенную команду cat.
// Выводит содержимое файлов или стандартного ввода.
type CatCommand struct{}

// NewCatCommand создает новый экземпляр команды cat.
func NewCatCommand() *CatCommand {
	return &CatCommand{}
}

// Name возвращает имя команды cat.
func (c *CatCommand) Name() string {
	return CatCommandName
}

// Execute выполняет команду cat.
// Если аргументы не переданы, читает из стандартного ввода.
// Иначе читает и выводит содержимое указанных файлов.
func (c *CatCommand) Execute(args []string, _ map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		_, err := io.Copy(stdout, stdin)
		if err != nil {
			fmt.Fprintf(stderr, "cat: error reading from stdin: %v\n", err)
			return 1
		}
		return 0
	}

	for _, filename := range args {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(stderr, "cat: %s: %v\n", filename, err)
			return 1
		}

		_, err = io.Copy(stdout, file)
		file.Close()

		if err != nil {
			fmt.Fprintf(stderr, "cat: %s: %v\n", filename, err)
			return 1
		}
	}

	return 0
}
