package builtins

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const WcCommandName = "wc"

// WcCommand реализует встроенную команду wc.
// Подсчитывает количество строк, слов и байтов в файле или стандартном вводе.
type WcCommand struct{}

// NewWcCommand создает новый экземпляр команды wc.
func NewWcCommand() *WcCommand {
	return &WcCommand{}
}

// Name возвращает имя команды wc.
func (w *WcCommand) Name() string {
	return WcCommandName
}

// Execute выполняет команду wc.
// Если аргументы не переданы, читает из стандартного ввода.
// Иначе читает указанный файл и подсчитывает статистику.
func (w *WcCommand) Execute(args []string, _ map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	var input io.Reader
	var filename string

	if len(args) == 0 {
		input = stdin
	} else {
		filename = args[0]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(stderr, "wc: %s: %v\n", filename, err)
			return 1
		}
		defer file.Close()
		input = file
	}

	lines, words, bytes := w.count(input)

	fmt.Fprintf(stdout, "%d %d %d", lines, words, bytes)
	if filename != "" {
		fmt.Fprintf(stdout, " %s", filename)
	}
	fmt.Fprintln(stdout)

	return 0
}

// count подсчитывает количество строк, слов и байтов в потоке ввода.
// Возвращает количество строк, слов и байтов соответственно.
func (w *WcCommand) count(input io.Reader) (int, int, int) {
	scanner := bufio.NewScanner(input)
	lines := 0
	words := 0
	bytes := 0

	for scanner.Scan() {
		line := scanner.Text()
		lines++
		bytes += len(line) + 1

		fields := strings.Fields(line)
		words += len(fields)
	}

	return lines, words, bytes
}
