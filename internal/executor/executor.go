package executor

import (
	"fmt"
	"os"
	osexec "os/exec"

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
func (exec *Executor) Execute(node parser.Node) error {
	switch n := node.(type) {
	case *parser.Command:
		return exec.executeCommand(n)
	case *parser.Pipeline:
		return exec.executePipeline(n)
	default:
		return fmt.Errorf("unknown node type: %T", node)
	}
}

// executeCommand выполняет отдельную команду.
// Устанавливает переменные окружения, определяет тип команды (встроенная/внешняя)
// и вызывает соответствующий метод выполнения.
func (exec *Executor) executeCommand(cmd *parser.Command) error {
	// Сохраняем состояние переменных для восстановления после выполнения
	type varState struct {
		wasLocal  bool
		oldValue  string
		wasGlobal bool
	}
	savedVars := make(map[string]varState)

	// Устанавливаем переменные окружения в локальном контексте
	for _, assignment := range cmd.Assignments {
		// Проверяем, существовала ли переменная в локальном окружении
		if oldValue, exists := exec.environment.GetLocal(assignment.Name); exists {
			savedVars[assignment.Name] = varState{
				wasLocal:  true,
				oldValue:  oldValue,
				wasGlobal: false,
			}
		} else if exec.environment.HasGlobal(assignment.Name) {
			// Переменная была только в глобальном окружении
			savedVars[assignment.Name] = varState{
				wasLocal:  false,
				oldValue:  "",
				wasGlobal: true,
			}
		} else {
			// Переменная не существовала
			savedVars[assignment.Name] = varState{
				wasLocal:  false,
				oldValue:  "",
				wasGlobal: false,
			}
		}
		exec.environment.Set(assignment.Name, assignment.Value.Value)
	}

	// Восстанавливаем переменные после выполнения команды
	defer func() {
		for name, state := range savedVars {
			if state.wasLocal {
				// Восстанавливаем старое локальное значение
				exec.environment.Set(name, state.oldValue)
			} else if state.wasGlobal {
				// Удаляем локальную переменную, чтобы вернуться к глобальной
				exec.environment.Unset(name)
			} else {
				// Переменная не существовала, удаляем её
				exec.environment.Unset(name)
			}
		}
	}()

	args := make([]string, len(cmd.Args))
	for i, arg := range cmd.Args {
		args[i] = arg.Value
	}

	if builtin, exists := exec.registry.Get(cmd.Name); exists {
		return exec.executeBuiltin(builtin, args)
	}

	return exec.executeExternal(cmd.Name, args)
}

// executePipeline выполняет пайплайн команд.
// В Phase 1 поддерживается только одна команда в пайплайне.
// В Phase 2 будет реализована полная поддержка пайплайнов.
func (exec *Executor) executePipeline(pipeline *parser.Pipeline) error {
	if len(pipeline.Commands) > 1 {
		return fmt.Errorf("pipeline execution not implemented in Phase 1")
	}

	if len(pipeline.Commands) == 0 {
		return fmt.Errorf("empty pipeline")
	}

	return exec.executeCommand(pipeline.Commands[0])
}

// executeBuiltin выполняет встроенную команду.
// Создает IO структуру и вызывает метод Execute встроенной команды.
// Передает переменные окружения во встроенную команду.
// Возвращает ошибку, если команда завершилась с ненулевым кодом возврата.
//
// Примечание: команда exit вызывает os.Exit(), который завершает процесс,
// поэтому код после вызова Execute для exit никогда не выполняется.
func (exec *Executor) executeBuiltin(builtin builtins.Builtin, args []string) error {
	io := builtins.NewIO()
	env := exec.environment.GetAllMap()

	exitCode := builtin.Execute(args, env, io.Stdin, io.Stdout, io.Stderr)

	// Если команда exit была вызвана, os.Exit() уже завершил процесс,
	// поэтому этот код не выполнится для exit
	if exitCode != 0 {
		return fmt.Errorf("command %s exited with code %d", builtin.Name(), exitCode)
	}

	return nil
}

// executeExternal выполняет внешнюю программу.
// Использует os/exec для запуска внешней команды с переданными аргументами.
// Передает переменные окружения и подключает стандартные потоки ввода/вывода.
func (exec *Executor) executeExternal(name string, args []string) error {
	cmd := osexec.Command(name, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = exec.environment.GetAll()

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("external command failed: %w", err)
	}

	return nil
}

func (exec *Executor) IsBuiltin(name string) bool {
	return exec.registry.IsBuiltin(name)
}

func (exec *Executor) ListBuiltins() []string {
	return exec.registry.List()
}

// SetEnvironment устанавливает окружение для исполнителя.
// Позволяет использовать общее окружение между Shell и Executor.
func (exec *Executor) SetEnvironment(env *environment.Environment) {
	exec.environment = env
}
