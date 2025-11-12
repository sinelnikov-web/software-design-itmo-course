package parser

// Node представляет узел абстрактного синтаксического дерева (AST).
// Все узлы AST должны реализовывать этот интерфейс.
type Node interface {
	String() string // Возвращает строковое представление узла
	Type() NodeType // Возвращает тип узла для классификации
}

// NodeType определяет тип узла в AST.
// Используется для классификации различных элементов синтаксического дерева.
type NodeType int

const (
	CommandNode    NodeType = iota // Узел команды
	PipelineNode                   // Узел пайплайна
	AssignmentNode                 // Узел присваивания переменной
	ArgumentNode                   // Узел аргумента
)

// Command представляет команду в AST.
// Содержит имя команды, аргументы и присваивания переменных окружения.
type Command struct {
	Name        string        // Имя команды (например, "echo", "cat")
	Args        []*Argument   // Аргументы команды
	Assignments []*Assignment // Присваивания переменных окружения
}

// Type возвращает тип узла Command.
func (c *Command) Type() NodeType {
	return CommandNode
}

// String возвращает строковое представление команды.
// Формат: "command arg1 arg2 ..."
func (c *Command) String() string {
	result := c.Name
	for _, arg := range c.Args {
		result += " " + arg.String()
	}
	return result
}

// Pipeline представляет пайплайн команд в AST.
// Содержит последовательность команд, соединенных оператором |.
type Pipeline struct {
	Commands []*Command // Список команд в пайплайне
}

// Type возвращает тип узла Pipeline.
func (p *Pipeline) Type() NodeType {
	return PipelineNode
}

// String возвращает строковое представление пайплайна.
// Формат: "cmd1 | cmd2 | cmd3"
func (p *Pipeline) String() string {
	if len(p.Commands) == 0 {
		return ""
	}

	result := p.Commands[0].String()
	for i := 1; i < len(p.Commands); i++ {
		result += " | " + p.Commands[i].String()
	}
	return result
}

// Assignment представляет присваивание переменной окружения в AST.
// Содержит имя переменной и её значение.
type Assignment struct {
	Name  string    // Имя переменной
	Value *Argument // Значение переменной
}

// Type возвращает тип узла Assignment.
func (a *Assignment) Type() NodeType {
	return AssignmentNode
}

// String возвращает строковое представление присваивания.
// Формат: "NAME=value"
func (a *Assignment) String() string {
	return a.Name + "=" + a.Value.String()
}

// QuoteType определяет тип кавычек для аргумента.
type QuoteType int

const (
	NoQuote     QuoteType = iota // Аргумент без кавычек
	SingleQuote                  // Одинарные кавычки (')
	DoubleQuote                  // Двойные кавычки (")
)

// Argument представляет аргумент команды в AST.
// Содержит значение аргумента и информацию о типе кавычек.
type Argument struct {
	Value     string    // Значение аргумента
	Quoted    bool      // Флаг: был ли аргумент в кавычках (для обратной совместимости)
	QuoteType QuoteType // Тип кавычек: NoQuote, SingleQuote или DoubleQuote
}

// Type возвращает тип узла Argument.
func (a *Argument) Type() NodeType {
	return ArgumentNode
}

// String возвращает строковое представление аргумента.
// Если аргумент был в кавычках, возвращает его в двойных кавычках.
func (a *Argument) String() string {
	if a.Quoted {
		return `"` + a.Value + `"`
	}
	return a.Value
}
