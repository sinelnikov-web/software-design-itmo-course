package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gocli/internal/environment"
	"gocli/internal/executor"
	"gocli/internal/lexer"
	"gocli/internal/parser"
)

// Shell представляет основную структуру командной оболочки.
// Содержит все необходимые компоненты для обработки пользовательского ввода:
// лексер для токенизации, парсер для построения AST и исполнитель для выполнения команд.
type Shell struct {
	executor    *executor.Executor       // Исполнитель команд (встроенные и внешние)
	lexer       *lexer.Lexer             // Лексер для разбора командной строки на токены
	parser      *parser.Parser           // Парсер для построения абстрактного синтаксического дерева
	environment *environment.Environment // Управление переменными окружения
}

// NewShell создает и инициализирует новый экземпляр командной оболочки.
// Возвращает готовую к использованию структуру Shell с настроенными компонентами.
func NewShell() *Shell {
	exec := executor.NewExecutor()
	env := environment.NewEnvironment()
	exec.SetEnvironment(env)

	return &Shell{
		executor:    exec,
		lexer:       lexer.NewLexer(),
		parser:      parser.NewParser(),
		environment: env,
	}
}

// Run запускает основной цикл командной оболочки (Read-Eval-Print Loop).
// Читает пользовательский ввод, обрабатывает команды и выводит результаты.
// Возвращает ошибку при завершении работы или при критических ошибках.
func (s *Shell) Run() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if err := s.processCommand(line); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	return scanner.Err()
}

func (s *Shell) processCommand(line string) error {
	tokens, err := s.lexer.Tokenize(line)
	if err != nil {
		return fmt.Errorf("lexical analysis failed: %w", err)
	}

	ast, err := s.parser.Parse(tokens)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	return s.executor.Execute(ast)
}
