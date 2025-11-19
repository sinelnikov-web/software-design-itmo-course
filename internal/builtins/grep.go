package builtins

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

const GrepCommandName = "grep"

// GrepCommand реализует встроенную команду grep.
// Выполняет поиск по регулярным выражениям в файлах или стандартном вводе.
type GrepCommand struct{}

// NewGrepCommand создает новый экземпляр команды grep.
func NewGrepCommand() *GrepCommand {
	return &GrepCommand{}
}

// Name возвращает имя команды grep.
func (g *GrepCommand) Name() string {
	return GrepCommandName
}

// Execute выполняет команду grep.
// Поддерживает флаги: -w (слово целиком), -i (регистронезависимый поиск), -A (строки после совпадения).
func (g *GrepCommand) Execute(args []string, _ map[string]string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	// Парсинг флагов с использованием стандартной библиотеки flag
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)
	fs.SetOutput(stderr) // Перенаправляем вывод ошибок flag в stderr

	wordBoundary := fs.Bool("w", false, "match whole word only")
	caseInsensitive := fs.Bool("i", false, "case-insensitive matching")
	afterContext := fs.Int("A", 0, "print N lines after match")

	// Парсим аргументы, пропуская первый (имя команды)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		fmt.Fprintf(stderr, "grep: %v\n", err)
		return 2
	}

	// Получаем оставшиеся аргументы после парсинга флагов
	remainingArgs := fs.Args()

	if len(remainingArgs) == 0 {
		fmt.Fprintf(stderr, "grep: pattern required\n")
		return 2
	}

	pattern := remainingArgs[0]
	filenames := remainingArgs[1:]

	// Строим регулярное выражение
	regexPattern := pattern
	if *wordBoundary {
		// Для поиска слова целиком добавляем границы слова
		// Используем \b для границ слова (non-word constituent characters)
		// Экранируем специальные символы в паттерне, чтобы он трактовался как литерал
		regexPattern = `\b` + regexp.QuoteMeta(pattern) + `\b`
	}
	// Если не используется -w, pattern используется как регулярное выражение напрямую

	// Компилируем регулярное выражение
	var re *regexp.Regexp
	var err error
	if *caseInsensitive {
		re, err = regexp.Compile("(?i)" + regexPattern)
	} else {
		re, err = regexp.Compile(regexPattern)
	}

	if err != nil {
		fmt.Fprintf(stderr, "grep: invalid regular expression: %v\n", err)
		return 2
	}

	// Если файлы не указаны, читаем из stdin
	if len(filenames) == 0 {
		return g.searchInReader(stdin, re, *afterContext, stdout, stderr, "")
	}

	// Обрабатываем каждый файл
	exitCode := 0
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(stderr, "grep: %s: %v\n", filename, err)
			exitCode = 1
			continue
		}

		fileExitCode := g.searchInReader(file, re, *afterContext, stdout, stderr, filename)
		file.Close()

		if fileExitCode != 0 {
			exitCode = fileExitCode
		}
	}

	return exitCode
}

// searchInReader выполняет поиск в потоке ввода.
// Возвращает код возврата: 0 если найдено совпадение, 1 если не найдено, 2 при ошибке.
func (g *GrepCommand) searchInReader(reader io.Reader, re *regexp.Regexp, afterContext int, stdout io.Writer, stderr io.Writer, filename string) int {
	scanner := bufio.NewScanner(reader)
	found := false
	linesAfter := 0 // Счетчик оставшихся строк для печати после совпадения

	for scanner.Scan() {
		line := scanner.Text()
		matched := re.MatchString(line)

		if matched {
			found = true
			// Печатаем текущую строку с совпадением
			if filename != "" {
				fmt.Fprintf(stdout, "%s:%s\n", filename, line)
			} else {
				fmt.Fprintln(stdout, line)
			}
			// Устанавливаем счетчик строк после совпадения
			linesAfter = afterContext
		} else if linesAfter > 0 {
			// Печатаем строку после совпадения с префиксом "-"
			if filename != "" {
				fmt.Fprintf(stdout, "%s-%s\n", filename, line)
			} else {
				fmt.Fprintf(stdout, "-%s\n", line)
			}
			linesAfter--
		}
		// Если linesAfter == 0 и нет совпадения, строка не печатается
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(stderr, "grep: error reading input: %v\n", err)
		return 2
	}

	if found {
		return 0
	}
	return 1
}
