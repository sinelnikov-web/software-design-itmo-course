package executor

import (
	"fmt"
	"os"
	"os/exec"

	"gocli/internal/builtins"
	"gocli/internal/environment"
	"gocli/internal/parser"
)

// Executor выполняет команды, представленные в виде AST.
// Различает встроенные команды и внешние программы, управляет их выполнением.
type Executor struct {
	registry    *builtins.Registry       // Реестр встроенных команд
	environment *environment.Environment // Управление переменными окружения
}

// NewExecutor создает новый экземпляр исполнителя.
// Инициализирует реестр встроенных команд и возвращает готовую структуру.
func NewExecutor() *Executor {
	return &Executor{
		registry:    builtins.NewRegistry(),
		environment: environment.NewEnvironment(),
	}
}

// Execute выполняет узел AST (команду или пайплайн).
// Определяет тип узла и вызывает соответствующий метод выполнения.
// Возвращает ошибку при неудачном выполнении команды.
func (e *Executor) Execute(node parser.Node) error {
	switch n := node.(type) {
	case *parser.Command:
		return e.executeCommand(n)
	case *parser.Pipeline:
		return e.executePipeline(n)
	default:
		return fmt.Errorf("unknown node type: %T", node)
	}
}

// executeCommand выполняет отдельную команду.
// Устанавливает переменные окружения, определяет тип команды (встроенная/внешняя)
// и вызывает соответствующий метод выполнения.
func (e *Executor) executeCommand(cmd *parser.Command) error {
	// Устанавливаем переменные окружения в локальном контексте
	for _, assignment := range cmd.Assignments {
		e.environment.Set(assignment.Name, assignment.Value.Value)
	}

	args := make([]string, len(cmd.Args))
	for i, arg := range cmd.Args {
		args[i] = arg.Value
	}

	if builtin, exists := e.registry.Get(cmd.Name); exists {
		return e.executeBuiltin(builtin, args)
	}

	return e.executeExternal(cmd.Name, args)
}

// executePipeline выполняет пайплайн команд.
// В Phase 1 поддерживается только одна команда в пайплайне.
// В Phase 2 будет реализована полная поддержка пайплайнов.
func (e *Executor) executePipeline(pipeline *parser.Pipeline) error {
	if len(pipeline.Commands) > 1 {
		return fmt.Errorf("pipeline execution not implemented in Phase 1")
	}

	if len(pipeline.Commands) == 0 {
		return fmt.Errorf("empty pipeline")
	}

	return e.executeCommand(pipeline.Commands[0])
}

// executeBuiltin выполняет встроенную команду.
// Создает IO структуру и вызывает метод Execute встроенной команды.
// Возвращает ошибку, если команда завершилась с ненулевым кодом возврата.
func (e *Executor) executeBuiltin(builtin builtins.Builtin, args []string) error {
	io := builtins.NewIO()

	exitCode := builtin.Execute(args, nil, io.Stdin, io.Stdout, io.Stderr)

	if exitCode != 0 {
		return fmt.Errorf("command %s exited with code %d", builtin.Name(), exitCode)
	}

	return nil
}

// executeExternal выполняет внешнюю программу.
// Использует os/exec для запуска внешней команды с переданными аргументами.
// Передает переменные окружения и подключает стандартные потоки ввода/вывода.
func (e *Executor) executeExternal(name string, args []string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = e.environment.GetAll()

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("external command failed: %w", err)
	}

	return nil
}

func (e *Executor) IsBuiltin(name string) bool {
	return e.registry.IsBuiltin(name)
}

func (e *Executor) ListBuiltins() []string {
	return e.registry.List()
}

// SetEnvironment устанавливает окружение для исполнителя.
// Позволяет использовать общее окружение между Shell и Executor.
func (e *Executor) SetEnvironment(env *environment.Environment) {
	e.environment = env
}
