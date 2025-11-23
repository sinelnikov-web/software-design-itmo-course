package expander

import (
	"fmt"
	"strings"
	"unicode"

	"gocli/internal/environment"
	"gocli/internal/parser"
)

// Expander выполняет подстановку переменных и обработку кавычек.
// Преобразует AST с переменными в окончательные аргументы команд.
type Expander struct {
	environment *environment.Environment
}

// NewExpander создает новый экземпляр expander.
func NewExpander(env *environment.Environment) *Expander {
	return &Expander{
		environment: env,
	}
}

// Expand выполняет подстановку переменных в AST.
// Обрабатывает подстановки $VAR в аргументах команд и присваиваниях.
// Возвращает расширенный AST и ошибку при некорректных подстановках.
func (e *Expander) Expand(node parser.Node) (parser.Node, error) {
	switch n := node.(type) {
	case *parser.Command:
		return e.expandCommand(n)
	case *parser.Pipeline:
		return e.expandPipeline(n)
	default:
		return nil, fmt.Errorf("unknown node type: %T", node)
	}
}

// expandCommand выполняет подстановку переменных в команде.
// Сначала расширяет assignments и устанавливает их в окружение временно,
// затем расширяет имя команды и аргументы (которые могут использовать эти переменные).
// Переменные остаются в окружении для Executor, но в пайплайнах они изолированы.
func (e *Expander) expandCommand(cmd *parser.Command) (*parser.Command, error) {
	expandedCmd := &parser.Command{
		Name:        cmd.Name,
		Args:        make([]*parser.Argument, 0, len(cmd.Args)),
		Assignments: make([]*parser.Assignment, 0, len(cmd.Assignments)),
	}

	// Сохраняем состояние переменных для восстановления после расширения
	savedVars := make(map[string]string)
	for _, assignment := range cmd.Assignments {
		if value, exists := e.environment.Get(assignment.Name); exists {
			savedVars[assignment.Name] = value
		}
	}

	// Сначала расширяем присваивания и устанавливаем их в окружение
	// Это позволяет использовать переменные из assignments в аргументах команды
	for _, assignment := range cmd.Assignments {
		expandedValue, err := e.expandArgument(assignment.Value)
		if err != nil {
			e.restoreVars(savedVars)
			return nil, fmt.Errorf("failed to expand assignment %s: %w", assignment.Name, err)
		}
		e.environment.Set(assignment.Name, expandedValue.Value)
		expandedCmd.Assignments = append(expandedCmd.Assignments, &parser.Assignment{
			Name:  assignment.Name,
			Value: expandedValue,
		})
	}

	// Расширяем имя команды и аргументы
	expandedName, err := e.expandString(cmd.Name)
	if err != nil {
		e.restoreVars(savedVars)
		return nil, fmt.Errorf("failed to expand command name: %w", err)
	}
	expandedCmd.Name = expandedName

	for _, arg := range cmd.Args {
		expandedArg, err := e.expandArgument(arg)
		if err != nil {
			e.restoreVars(savedVars)
			return nil, fmt.Errorf("failed to expand argument: %w", err)
		}
		expandedCmd.Args = append(expandedCmd.Args, expandedArg)
	}

	// Переменные остаются в окружении для Executor
	// В пайплайнах они будут откатены в expandPipeline
	return expandedCmd, nil
}

// restoreVars восстанавливает сохраненные переменные в окружении.
func (e *Expander) restoreVars(savedVars map[string]string) {
	for name, value := range savedVars {
		e.environment.Set(name, value)
	}
}

// expandPipeline выполняет подстановку переменных в пайплайне.
// Переменные каждой команды изолированы - после расширения команды они откатываются.
func (e *Expander) expandPipeline(pipeline *parser.Pipeline) (*parser.Pipeline, error) {
	expandedPipeline := &parser.Pipeline{
		Commands: make([]*parser.Command, 0, len(pipeline.Commands)),
	}

	for _, cmd := range pipeline.Commands {
		// Сохраняем состояние переменных перед расширением команды
		savedVars := make(map[string]string)
		for _, assignment := range cmd.Assignments {
			if value, exists := e.environment.Get(assignment.Name); exists {
				savedVars[assignment.Name] = value
			}
		}

		expandedCmd, err := e.expandCommand(cmd)
		if err != nil {
			return nil, err
		}

		// Откатываем переменные этой команды, чтобы они не влияли на следующие команды
		// Executor установит их снова при выполнении
		for _, assignment := range expandedCmd.Assignments {
			if _, wasSaved := savedVars[assignment.Name]; wasSaved {
				e.environment.Set(assignment.Name, savedVars[assignment.Name])
			} else {
				e.environment.Unset(assignment.Name)
			}
		}

		expandedPipeline.Commands = append(expandedPipeline.Commands, expandedCmd)
	}

	return expandedPipeline, nil
}

// expandArgument выполняет подстановку переменных в аргументе.
// Учитывает тип кавычек: одинарные кавычки не расширяются, двойные - расширяются.
func (e *Expander) expandArgument(arg *parser.Argument) (*parser.Argument, error) {
	// Если аргумент был в одинарных кавычках, подстановки не выполняются
	if arg.QuoteType == parser.SingleQuote {
		return &parser.Argument{
			Value:     arg.Value,
			Quoted:    false, // После обработки кавычки убираются
			QuoteType: parser.NoQuote,
		}, nil
	}

	// Для двойных кавычек и некавыченных аргументов выполняем подстановки
	expanded, err := e.expandString(arg.Value)
	if err != nil {
		return nil, err
	}
	return &parser.Argument{
		Value:     expanded,
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}, nil
}

// expandString выполняет подстановку переменных в строке.
// Одинарные кавычки обрабатываются на уровне expandArgument, здесь всегда выполняются подстановки.
// Обрабатывает экранирование обратными слешами и поддерживает синтаксис $VAR и ${VAR}.
func (e *Expander) expandString(s string) (string, error) {
	var result strings.Builder
	var i int

	for i < len(s) {
		dollarIdx := e.findNextDollar(s, i)
		if dollarIdx == -1 {
			// $ не найден, записываем остаток и выходим
			result.WriteString(s[i:])
			break
		}

		backslashCount := e.countBackslashesBefore(s, dollarIdx)
		if e.handleEscapedDollar(&result, s, &i, dollarIdx, backslashCount) {
			continue
		}

		e.writeTextBeforeDollar(&result, s, i, dollarIdx, backslashCount)

		newPos, err := e.expandVariable(&result, s, dollarIdx)
		if err != nil {
			return "", err
		}
		i = newPos
	}

	return result.String(), nil
}

// findNextDollar ищет следующий символ $ в строке начиная с позиции start.
// Возвращает индекс найденного $ или -1, если $ не найден.
func (e *Expander) findNextDollar(s string, start int) int {
	for j := start; j < len(s); j++ {
		if s[j] == '$' {
			return j
		}
	}
	return -1
}

// countBackslashesBefore подсчитывает количество последовательных обратных слешей
// перед позицией dollarIdx в строке s.
func (e *Expander) countBackslashesBefore(s string, dollarIdx int) int {
	count := 0
	for k := dollarIdx - 1; k >= 0 && s[k] == '\\'; k-- {
		count++
	}
	return count
}

// handleEscapedDollar обрабатывает случай, когда $ экранирован обратным слешем.
// Если $ экранирован (нечетное количество \), записывает текст до $ без последнего \
// и сам $ как обычный символ. Возвращает true, если $ был экранирован.
func (e *Expander) handleEscapedDollar(result *strings.Builder, s string, i *int, dollarIdx, backslashCount int) bool {
	if backslashCount%2 == 1 {
		// $ экранирован - записываем текст до $ без последнего экранирующего \
		result.WriteString(s[*i : dollarIdx-1])
		result.WriteRune('$')
		*i = dollarIdx + 1
		return true
	}
	return false
}

// writeTextBeforeDollar записывает текст до символа $ в результат.
// Обрабатывает пары обратных слешей: каждая пара \\ становится одним \ в результате.
func (e *Expander) writeTextBeforeDollar(result *strings.Builder, s string, start, dollarIdx, backslashCount int) {
	if backslashCount > 0 {
		// Записываем текст до начала последовательности \
		result.WriteString(s[start : dollarIdx-backslashCount])
		// Записываем половину обратных слешей (каждая пара \\ → один \)
		for j := 0; j < backslashCount/2; j++ {
			result.WriteRune('\\')
		}
	} else {
		// Нет обратных слешей - записываем весь текст до $
		result.WriteString(s[start:dollarIdx])
	}
}

// expandVariable обрабатывает подстановку переменной после символа $.
// Поддерживает синтаксис ${VAR} и $VAR с умным fallback для $VAR_suffix.
// Возвращает новую позицию в строке после обработки переменной.
func (e *Expander) expandVariable(result *strings.Builder, s string, dollarIdx int) (int, error) {
	// $ в конце строки - оставляем как есть
	if dollarIdx+1 >= len(s) {
		result.WriteRune('$')
		return len(s), nil
	}

	next := s[dollarIdx+1]
	if next == '{' {
		// Подстановка вида ${VAR}
		return e.expandBracedVariable(result, s, dollarIdx)
	}
	if isValidVariableStart(rune(next)) {
		// Подстановка вида $VAR
		return e.expandSimpleVariable(result, s, dollarIdx), nil
	}
	// $ не является началом переменной
	result.WriteRune('$')
	return dollarIdx + 1, nil
}

// expandBracedVariable обрабатывает подстановку переменной в фигурных скобках ${VAR}.
// Возвращает новую позицию в строке после закрывающей скобки.
func (e *Expander) expandBracedVariable(result *strings.Builder, s string, dollarIdx int) (int, error) {
	closeIdx := -1
	for j := dollarIdx + 2; j < len(s); j++ {
		if s[j] == '}' {
			closeIdx = j
			break
		}
	}
	if closeIdx == -1 {
		return 0, fmt.Errorf("unclosed ${ variable")
	}
	varName := s[dollarIdx+2 : closeIdx]
	result.WriteString(e.getVariableValue(varName))
	return closeIdx + 1, nil
}

// expandSimpleVariable обрабатывает подстановку переменной без фигурных скобок $VAR.
// Использует умный fallback: если переменная $VAR_suffix не найдена,
// пробует более короткие префиксы (например, $VAR).
// Возвращает новую позицию в строке после имени переменной.
func (e *Expander) expandSimpleVariable(result *strings.Builder, s string, dollarIdx int) int {
	startIdx := dollarIdx + 1
	endIdx := startIdx + 1
	for endIdx < len(s) && isValidVariableChar(rune(s[endIdx])) {
		endIdx++
	}

	varName := s[startIdx:endIdx]
	value := e.getVariableValue(varName)

	// Fallback: если переменная не найдена, пробуем более короткие префиксы
	// Например, для $VAR_suffix сначала ищем VAR_suffix, затем VAR
	if value == "" && endIdx > startIdx+1 {
		for k := endIdx - 1; k > startIdx; k-- {
			candidate := s[startIdx:k]
			if v := e.getVariableValue(candidate); v != "" {
				value = v
				endIdx = k
				break
			}
		}
	}

	result.WriteString(value)
	return endIdx
}

// isValidVariableStart проверяет, может ли символ быть началом имени переменной.
func isValidVariableStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// isValidVariableChar проверяет, может ли символ быть частью имени переменной.
func isValidVariableChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// getVariableValue получает значение переменной из окружения.
// Если переменная не существует, возвращает пустую строку (как в POSIX).
func (e *Expander) getVariableValue(name string) string {
	if v, ok := e.environment.Get(name); ok {
		return v
	}
	return ""
}
