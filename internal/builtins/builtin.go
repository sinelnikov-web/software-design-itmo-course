package builtins

import (
	"io"
	"os"
)

// Builtin определяет интерфейс для встроенных команд.
// Все встроенные команды должны реализовывать этот интерфейс.
type Builtin interface {
	// Execute выполняет команду с переданными аргументами и потоками ввода/вывода.
	// Возвращает код возврата (0 для успеха, ненулевое значение для ошибки).
	Execute(args []string, env map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int
	// Name возвращает имя команды.
	Name() string
}

// IO представляет стандартные потоки ввода/вывода.
// Используется для передачи потоков встроенным командам.
type IO struct {
	Stdin  io.Reader // Стандартный ввод
	Stdout io.Writer // Стандартный вывод
	Stderr io.Writer // Стандартный поток ошибок
}

// NewIO создает новую структуру IO с системными потоками ввода/вывода.
// Возвращает готовую к использованию структуру с подключенными потоками.
func NewIO() *IO {
	return &IO{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}
