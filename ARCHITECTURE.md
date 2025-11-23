# software-design-itmo-course

## Команда

1. Синельников Максим
2. Малдыбаев Руслан Ашимович
3. Короткевич Леонид Витальевич
4. Батыров Амир Фанисович

## 1. Цели и область

Спроектировать минималистичный, расширяемый интерпретатор командной строки (CLI Shell), поддерживающий:

* **Команды (builtins):** `cat`, `echo`, `wc`, `pwd`, `exit`.
* **Кавычки:** одинарные (full) и двойные (weak) quoting.
* **Окружение:** присваивания `name=value`, подстановки через `$NAME`.
* **Вызов внешних программ**, если команда не является встроенной.
* **Пайплайны** через оператор `|` (конвейеры).

Нефункциональные требования:

* Простота добавления новых команд.
* Чёткое разделение ответственности между компонентами.
* Компонентная структура (не "клубок классов").
* Наличие текстового описания архитектуры и структурной диаграммы.
* Отражение **этапности** реализации (Фаза 1 → Фаза 2).

---

## 2. Высокоуровневая идея

Интерпретатор разделён на независимые слои:

1. **REPL** (ввод/луп): читает строку, отдаёт в пайплайн обработки, печатает результат/ошибки.
2. **Lexer/Tokenizer**: превращает исходную строку в токены (слова, кавычки, операторы, присваивания, пайпы).
3. **Parser**: строит **AST**/IR команды/конвейера (pipeline).
4. **Expander**: выполняет подстановки переменных и обработку кавычек, формируя окончательные аргументы.
5. **Executor**: исполняет AST: вызывает builtins или внешние программы, соединяет их потоками в пайплайне.
6. **Builtins Registry**: реестр встроенных команд с единым интерфейсом.
7. **Environment**: абстракция окружения (get/set).
8. **IO/Pipes**: абстракция стандартных потоков и каналов (stdin/stdout/stderr, pipe).
9. **Errors & Diagnostics**: единый формат ошибок.

Такое разбиение позволяет:

* Менять/тестировать каждый компонент отдельно.
* Легко добавлять команды через регистрацию в реестре.
* Позже подключить расширения (алиасы, редиректы `>`, `<`, `>>`, и т.д.).

---

## 3. Модель данных (IR/AST)

### 3.1 Узлы AST

* **`Command`**: { `name: string`, `args: string[]`, `assignments?: Record<string,string>` }
* **`Pipeline`**: { `stages: (Command)[]` }

Assignments хранят присваивания, стоящие **перед** командой в качестве переменных окружения

### 3.2 Токены

`WORD`, `SQUOTE`, `DQUOTE`, `ASSIGN` (имя=значение), `PIPE` (`|`).

---

## 4. Диаграммы (Mermaid)

### 4.1 Последовательность (пример `echo 123 | wc`)

```mermaid
sequenceDiagram
    participant U as User
    participant R as REPL
    participant L as Lexer
    participant P as Parser
    participant E as Expander
    participant X as Executor
    participant B as Builtins


    U->>R: echo 123 | wc\n
    R->>L: tokenize()
    L-->>R: [WORD echo, WORD 123, PIPE, WORD wc]
    R->>P: parse(tokens)
    P-->>R: Pipeline([Command(echo,[123]), Command(wc,[])])
    R->>E: expand(ast, env)
    E-->>R: (expanded substitutions)
    R->>X: execute(pipeline)
    X->>B: run(echo,[123]) -> stdout
    X->>B: run(wc,stdin=echoOut)
    B-->>R: prints "1 1 3"
```

---

### 4.2 Классы

```mermaid
classDiagram
    %% Главная точка входа
    class Main {
        +main()
    }

    %% REPL компонент
    class Shell {
        -executor: Executor
        -lexer: Lexer
        -parser: Parser
        -environment: Environment
        +NewShell() Shell
        +Run() error
        -processCommand(line string) error
    }

    %% Лексический анализатор
    class Lexer {
        +NewLexer() Lexer
        +Tokenize(input string) ([]Token, error)
        -isValidVariableName(name string) bool
    }

    class Token {
        +Type: TokenType
        +Value: string
        +String() string
    }

    class TokenType {
        <<enumeration>>
        WORD
        PIPE
        SQUOTE
        DQUOTE
        ASSIGN
    }

    %% Синтаксический анализатор
    class Parser {
        +NewParser() Parser
        +Parse(tokens []Token) (Node, error)
        -parsePipeline(tokens []Token) ([]*Command, error)
        -parseCommand(tokens []Token) (*Command, error)
        -createArgument(token Token) *Argument
    }

    %% AST узлы
    class Node {
        <<interface>>
        +String() string
        +Type() NodeType
    }

    class NodeType {
        <<enumeration>>
        CommandNode
        PipelineNode
        AssignmentNode
        ArgumentNode
    }

    class Command {
        +Name: string
        +Args: []*Argument
        +Assignments: []*Assignment
        +String() string
        +Type() NodeType
    }

    class Pipeline {
        +Commands: []*Command
        +String() string
        +Type() NodeType
    }

    class Assignment {
        +Name: string
        +Value: *Argument
        +String() string
        +Type() NodeType
    }

    class Argument {
        +Value: string
        +Quoted: bool
        +String() string
        +Type() NodeType
    }

    %% Исполнитель команд
    class Executor {
        -registry: Registry
        -environment: Environment
        +NewExecutor() Executor
        +Execute(node Node) error
        +IsBuiltin(name string) bool
        +ListBuiltins() []string
        +SetEnvironment(env Environment)
        -executeCommand(cmd *Command) error
        -executePipeline(pipeline *Pipeline) error
        -executeBuiltin(builtin Builtin, args []string) error
        -executeExternal(name string, args []string) error
    }

    %% Реестр встроенных команд
    class Registry {
        -commands: map[string]Builtin
        +NewRegistry() Registry
        +Register(command Builtin)
        +Get(name string) (Builtin, bool)
        +List() []string
        +IsBuiltin(name string) bool
        +String() string
    }

    %% Интерфейс встроенных команд
    class Builtin {
        <<interface>>
        +Execute(args []string, env map[string]string, stdin Reader, stdout Writer, stderr Writer) int
        +Name() string
    }

    class IO {
        +Stdin: Reader
        +Stdout: Writer
        +Stderr: Writer
        +NewIO() IO
    }

    %% Встроенные команды
    class CatCommand {
        +Name() string
        +Execute(...) int
        -processFile(filename string, stdout Writer, stderr Writer) int
    }

    class EchoCommand {
        +Name() string
        +Execute(...) int
    }

    class WcCommand {
        +Name() string
        +Execute(...) int
        -count(input Reader) (int, int, int)
    }

    class PwdCommand {
        +Name() string
        +Execute(...) int
    }

    class ExitCommand {
        +Name() string
        +Execute(...) int
    }

    %% Управление окружением
    class Environment {
        -global: map[string]string
        -local: map[string]string
        +NewEnvironment() Environment
        +Set(name string, value string)
        +Get(name string) (string, bool)
        +GetAll() []string
        +Unset(name string)
        +ClearLocal()
        +ListLocal() map[string]string
        +ListGlobal() map[string]string
    }

    %% Связи между компонентами
    Main --> Shell: создает
    
    Shell --> Lexer: использует
    Shell --> Parser: использует
    Shell --> Executor: использует
    Shell --> Environment: управляет
    
    Lexer --> Token: создает
    Token --> TokenType: использует
    
    Parser --> Token: принимает
    Parser --> Node: создает
    Parser --> Command: создает
    Parser --> Pipeline: создает
    Parser --> Argument: создает
    Parser --> Assignment: создает
    
    Node <|.. Command: реализует
    Node <|.. Pipeline: реализует
    Node <|.. Assignment: реализует
    Node <|.. Argument: реализует
    Command --> NodeType: использует
    Pipeline --> NodeType: использует
    Assignment --> NodeType: использует
    Argument --> NodeType: использует
    
    Pipeline --> Command: содержит
    Command --> Argument: содержит
    Command --> Assignment: содержит
    Assignment --> Argument: содержит
    
    Executor --> Registry: использует
    Executor --> Environment: использует
    Executor --> Node: выполняет
    Executor --> Command: выполняет
    Executor --> Pipeline: выполняет
    Executor --> Builtin: вызывает
    Executor --> IO: создает
    
    Registry --> Builtin: управляет
    Registry --> CatCommand: регистрирует
    Registry --> EchoCommand: регистрирует
    Registry --> WcCommand: регистрирует
    Registry --> PwdCommand: регистрирует
    Registry --> ExitCommand: регистрирует
    
    Builtin <|.. CatCommand: реализует
    Builtin <|.. EchoCommand: реализует
    Builtin <|.. WcCommand: реализует
    Builtin <|.. PwdCommand: реализует
    Builtin <|.. ExitCommand: реализует
    
    Environment ..> Main: системное окружение
```

## 5. Правила лексинга и парсинга

### 5.1 Лексер

* Разбивает по пробелам вне кавычек.
* Распознаёт `|` как отдельный токен `PIPE`.
* Одинарные кавычки `'...'` - содержимое не интерпретируется (full quoting).
* Двойные кавычки `"..."` - разрешены подстановки `$VAR` (weak quoting), backslash-escape опционально.
* `name=value` в начале слова или перед командой - токен `ASSIGN`.

### 5.2 Парсер

```
command   := { assignment } WORD { argument } ;
pipeline  := (command|pipeline) | (command|pipeline) | ... ;
assignment:= NAME '=' VALUE ;
argument  := WORD | SQUOTE | DQUOTE ;
```

Где `WORD` может включать фрагменты с `$VAR` (обрабатываются позже на фазе expand).

---

## 6. Подстановки и кавычки (Expander)

* **Переменные:** `$NAME` → подстановка из `Environment` или `LocalEnv`. Имя: `[A-Za-z_][A-Za-z0-9_]*`.
* **Одинарные кавычки `'...'`:** содержимое берётся как есть (никаких подстановок и экранирований).
* **Двойные кавычки `"..."`:** выполняются подстановки `$NAME`; сам символ `"` не входит в итоговый аргумент.
* **Конкатенация:** соседние лексемы склеиваются: `"$x"$y` → значение `x` + значение `y`.
* **Присваивания:** `VAR=VAL` без команды → изменить `LocalEnv`;

Ошибки подстановки (нет переменной) → подставить пустую строку (как в POSIX), либо configurable.

---

## 7. Исполнение (Executor)

### 7.1 Builtins vs внешние

* Проверка в `BuiltinsRegistry`.
* Если нет - запуск внешней программы через `ExternalProcessRunner` (обёртка над системным `exec`/`spawn`).

### 7.2 Пайплайны

* Подключить `stdout` i-й команды к `stdin` (i+1)-й.
* Запускать стадии последовательно.
* Дождаться завершения всех стадий, вернуть код последней стадии.

### 7.3 Контракты IO

* Все команды (встроенные и внешние) читают из `io.stdin`, пишут в `io.stdout`, ошибки в `io.stderr`.

---

## 8. Встроенные команды (Builtins)

Единый интерфейс `Builtin.execute(args, env, io): int`.

**Расширяемость:** Новая команда добавляется реализацией интерфейса `Builtin` и регистрацией `registry.register("name", impl)`.

## 9. Окружене (Environment)

Назначение: единый источник правды для переменных окружения Shell и временных присваиваний перед командами.

Ответственность и границы

* Хранение пар name → value (все значения — строки).
* Доступ к переменным для Expander/Executor.
* Скопы: различает глобальное окружение Shell и локальные временные присваивания для конкретного запуска команды/пайплайна.

Интерфейс

```
interface Environment {
 get(name: string): string | undefined
 set(name: string, value: string): void
 has(name: string): boolean
 all(): Map<string, string> 
}
```

---

## 10. Этапность реализации

### Фаза 1 (базовый Shell)

* Реализовать: REPL → Lexer → Parser → Executor → Builtins Registry → Environment (минимально, для внешних команд) → IO.
* Поддерживаемые фичи: одиночная команда **без** `$VAR`, **без** `|`. Кавычки допускаются, но трактуются как група слов (без expand логики).
* Builtins: `cat`, `echo`, `wc`, `pwd`, `exit`.
* Вызов внешних программ: да.

### Фаза 2 (подстановки и пайплайны)

* Включить `Expander` и полноценную логику кавычек/переменных.
* Поддержка `|` и конвейерного исполнения, обработка нескольких стадий.

**Замечание по коду:** Фаза 1 уже использует интерфейсы/контракты, чтобы в Фазе 2 подключить `Expander` и `Pipeline` без перелома архитектуры.

---

## 11. Примеры сценариев

```text
> echo "Hello, world!"
Hello, world!

> FILE=example.txt
> cat $FILE
Some example text

> cat example.txt | wc
1 3 18

> echo 123 | wc
1 1 3

> x=ex
> y=it
> $x$y   # expand → exit -> исполнение exit
```

---
