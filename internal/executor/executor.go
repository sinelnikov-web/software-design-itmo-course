package executor

import (
	"fmt"
	"io"
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

// varState хранит состояние переменной перед её изменением.
// Используется для восстановления переменных после выполнения команды.
type varState struct {
	wasLocal  bool   // Была ли переменная в локальном окружении
	oldValue  string // Старое значение (если была локальной)
	wasGlobal bool   // Была ли переменная только в глобальном окружении
}

// saveAndSetVariables сохраняет текущее состояние переменных и устанавливает новые значения.
// Возвращает map сохраненных состояний для последующего восстановления.
func (exec *Executor) saveAndSetVariables(assignments []*parser.Assignment) map[string]varState {
	savedVars := make(map[string]varState)

	for _, assignment := range assignments {
		// Сохраняем текущее состояние переменной
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
		// Устанавливаем новое значение
		exec.environment.Set(assignment.Name, assignment.Value.Value)
	}

	return savedVars
}

// restoreVariables восстанавливает сохраненное состояние переменных.
// Используется для отката временных переменных после выполнения команды.
func (exec *Executor) restoreVariables(savedVars map[string]varState) {
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
}

// executeCommand выполняет отдельную команду.
// Устанавливает переменные окружения, определяет тип команды (встроенная/внешняя)
// и вызывает соответствующий метод выполнения.
// Временные переменные из assignments автоматически восстанавливаются после выполнения.
func (exec *Executor) executeCommand(cmd *parser.Command) error {
	// Сохраняем и устанавливаем переменные окружения
	savedVars := exec.saveAndSetVariables(cmd.Assignments)

	// Восстанавливаем переменные после выполнения команды
	defer exec.restoreVariables(savedVars)

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
// Подключает stdout i-й команды к stdin (i+1)-й команды через pipe.
//
// Важно: команды выполняются параллельно (в отдельных goroutines), а не последовательно.
// Это необходимо для правильной работы pipe - команды должны читать и писать одновременно.
// Например, в пайплайне `echo hello | wc` команда `wc` должна начать читать данные
// сразу после того, как `echo` начнет их писать, иначе pipe заблокируется.
//
// Возвращается код последней команды (как в POSIX shell).
// Ошибки промежуточных команд не прерывают пайплайн, но сохраняются для диагностики.
func (exec *Executor) executePipeline(pipeline *parser.Pipeline) error {
	if len(pipeline.Commands) == 0 {
		return fmt.Errorf("empty pipeline")
	}

	// Если команда одна, выполняем её без пайплайна
	if len(pipeline.Commands) == 1 {
		return exec.executeCommand(pipeline.Commands[0])
	}

	// Создаем pipes между командами
	// Для N команд нужно N-1 pipe: между каждой парой соседних команд
	pipes := make([]*io.PipeWriter, len(pipeline.Commands)-1)
	readers := make([]io.Reader, len(pipeline.Commands))

	// Первая команда читает из os.Stdin
	readers[0] = os.Stdin

	// Создаем pipes для промежуточных команд
	// Каждый pipe соединяет stdout команды i с stdin команды i+1
	for i := 0; i < len(pipeline.Commands)-1; i++ {
		r, w := io.Pipe()
		pipes[i] = w
		readers[i+1] = r
	}

	// Запускаем все команды параллельно в goroutines
	// Это критично для работы pipe - команды должны работать одновременно
	errors := make(chan error, len(pipeline.Commands))

	for i, cmd := range pipeline.Commands {
		i := i // Захватываем переменную для замыкания в goroutine
		cmd := cmd

		go func() {
			// Определяем stdout для команды
			// Промежуточные команды пишут в pipe, последняя - в os.Stdout
			var stdout io.Writer = os.Stdout
			if i < len(pipeline.Commands)-1 {
				stdout = pipes[i]
			}

			// Выполняем команду с правильными потоками ввода/вывода
			err := exec.executeCommandInPipeline(cmd, readers[i], stdout, os.Stderr)

			// Закрываем pipe после записи (если это не последняя команда)
			// Это сигнализирует следующей команде, что данных больше не будет
			if i < len(pipeline.Commands)-1 {
				pipes[i].Close()
			}

			// Отправляем результат в канал ошибок
			if err != nil {
				errors <- fmt.Errorf("command %d (%s) failed: %w", i, cmd.Name, err)
			} else {
				errors <- nil
			}
		}()
	}

	// Ждем завершения всех команд и собираем ошибки
	// В POSIX shell код возврата пайплайна равен коду последней команды
	var lastErr error
	for i := 0; i < len(pipeline.Commands); i++ {
		if err := <-errors; err != nil {
			// Сохраняем ошибку последней команды - она определяет код возврата пайплайна
			if i == len(pipeline.Commands)-1 {
				lastErr = err
			}
			// Ошибки промежуточных команд не прерывают пайплайн
			// (в bash/zsh промежуточные ошибки игнорируются, если не установлен set -e)
		}
	}

	return lastErr
}

// executeCommandInPipeline выполняет команду в контексте пайплайна.
// stdin, stdout, stderr определяют потоки ввода/вывода для команды.
// Временные переменные из assignments автоматически восстанавливаются после выполнения.
func (exec *Executor) executeCommandInPipeline(
	cmd *parser.Command,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	// Сохраняем и устанавливаем переменные окружения
	savedVars := exec.saveAndSetVariables(cmd.Assignments)

	// Восстанавливаем переменные после выполнения команды
	defer exec.restoreVariables(savedVars)

	args := make([]string, len(cmd.Args))
	for j, arg := range cmd.Args {
		args[j] = arg.Value
	}

	// Выполняем команду
	if builtin, exists := exec.registry.Get(cmd.Name); exists {
		return exec.executeBuiltinInPipeline(builtin, args, stdin, stdout, stderr)
	}

	return exec.executeExternalInPipeline(cmd.Name, args, stdin, stdout, stderr)
}

// executeBuiltinInPipeline выполняет встроенную команду в контексте пайплайна.
func (exec *Executor) executeBuiltinInPipeline(
	builtin builtins.Builtin,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	env := exec.environment.GetAllMap()

	// Создаем IO структуру с переданными потоками
	io := &builtins.IO{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	exitCode := builtin.Execute(args, env, io.Stdin, io.Stdout, io.Stderr)

	// Команда exit в пайплайне завершает весь процесс
	// Это стандартное поведение для большинства shell
	if exitCode != 0 {
		return fmt.Errorf("command %s exited with code %d", builtin.Name(), exitCode)
	}

	return nil
}

// executeExternalInPipeline выполняет внешнюю программу в контексте пайплайна.
func (exec *Executor) executeExternalInPipeline(
	name string,
	args []string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	cmd := osexec.Command(name, args...)

	// Устанавливаем потоки
	if stdin != nil {
		cmd.Stdin = stdin
	} else {
		cmd.Stdin = os.Stdin
	}

	if stdout != nil {
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = os.Stdout
	}

	if stderr != nil {
		cmd.Stderr = stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	cmd.Env = exec.environment.GetAll()

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("external command failed: %w", err)
	}

	return nil
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
