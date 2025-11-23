package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// Lexer выполняет лексический анализ командной строки.
// Разбивает входную строку на токены, обрабатывая кавычки, пробелы и специальные символы.
type Lexer struct{}

// NewLexer создает новый экземпляр лексера.
// Возвращает готовую к использованию структуру Lexer.
func NewLexer() *Lexer {
	return &Lexer{}
}

// Tokenize выполняет лексический анализ входной строки.
// Разбирает строку на токены, учитывая кавычки, пробелы и специальные символы.
// Возвращает массив токенов и ошибку при некорректном вводе (например, незакрытые кавычки).
func (l *Lexer) Tokenize(input string) ([]Token, error) {
	state := &tokenizeState{
		tokens:        []Token{},
		current:       strings.Builder{},
		inSingleQuote: false,
		inDoubleQuote: false,
	}

	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		if err := l.processChar(char, state); err != nil {
			return nil, err
		}
	}

	return l.finalizeTokens(state)
}

// tokenizeState хранит состояние процесса токенизации.
type tokenizeState struct {
	tokens        []Token
	current       strings.Builder
	inSingleQuote bool
	inDoubleQuote bool
}

// processChar обрабатывает один символ в процессе токенизации.
func (l *Lexer) processChar(char rune, state *tokenizeState) error {
	switch {
	case char == '\'' && !state.inDoubleQuote:
		return l.handleSingleQuote(state)
	case char == '"' && !state.inSingleQuote:
		return l.handleDoubleQuote(state)
	case char == '|' && !state.inSingleQuote && !state.inDoubleQuote:
		l.handlePipe(state)
	case unicode.IsSpace(char) && !state.inSingleQuote && !state.inDoubleQuote:
		l.handleSpace(state)
	case char == '=' && !state.inSingleQuote && !state.inDoubleQuote:
		l.handleAssignment(char, state)
	default:
		state.current.WriteRune(char)
	}
	return nil
}

// handleSingleQuote обрабатывает одинарные кавычки.
func (l *Lexer) handleSingleQuote(state *tokenizeState) error {
	if state.inSingleQuote {
		// Закрытие одинарных кавычек: сохраняем содержимое как SQUOTE токен
		state.tokens = append(state.tokens, Token{Type: SQUOTE, Value: state.current.String()})
		state.current.Reset()
		state.inSingleQuote = false
	} else {
		// Открытие одинарных кавычек: сохраняем накопленное слово, если есть
		l.flushCurrentWord(state)
		state.inSingleQuote = true
	}
	return nil
}

// handleDoubleQuote обрабатывает двойные кавычки.
func (l *Lexer) handleDoubleQuote(state *tokenizeState) error {
	if state.inDoubleQuote {
		// Закрытие двойных кавычек: сохраняем содержимое как DQUOTE токен
		state.tokens = append(state.tokens, Token{Type: DQUOTE, Value: state.current.String()})
		state.current.Reset()
		state.inDoubleQuote = false
	} else {
		// Открытие двойных кавычек: сохраняем накопленное слово, если есть
		l.flushCurrentWord(state)
		state.inDoubleQuote = true
	}
	return nil
}

// handlePipe обрабатывает оператор пайплайна.
func (l *Lexer) handlePipe(state *tokenizeState) {
	l.flushCurrentWord(state)
	state.tokens = append(state.tokens, Token{Type: PIPE, Value: "|"})
}

// handleSpace обрабатывает пробельные символы.
func (l *Lexer) handleSpace(state *tokenizeState) {
	l.flushCurrentWord(state)
}

// handleAssignment обрабатывает оператор присваивания.
func (l *Lexer) handleAssignment(char rune, state *tokenizeState) {
	if state.current.Len() > 0 {
		word := state.current.String()
		// Если накопленная строка - валидное имя переменной, создаем ASSIGN токен
		if l.isValidVariableName(word) {
			state.tokens = append(state.tokens, Token{Type: ASSIGN, Value: word})
			state.current.Reset()
		} else {
			// Иначе добавляем '=' как часть слова
			state.current.WriteRune(char)
		}
	} else {
		// Если нет накопленного слова, '=' становится частью нового слова
		state.current.WriteRune(char)
	}
}

// flushCurrentWord сохраняет накопленное слово как WORD токен, если оно не пустое.
func (l *Lexer) flushCurrentWord(state *tokenizeState) {
	if state.current.Len() > 0 {
		state.tokens = append(state.tokens, Token{Type: WORD, Value: state.current.String()})
		state.current.Reset()
	}
}

// finalizeTokens завершает токенизацию и проверяет корректность.
func (l *Lexer) finalizeTokens(state *tokenizeState) ([]Token, error) {
	// Обработка оставшегося содержимого после завершения цикла
	if state.current.Len() > 0 {
		// Если осталось содержимое, но мы все еще в кавычках - ошибка
		if state.inSingleQuote || state.inDoubleQuote {
			return nil, fmt.Errorf("unclosed quote")
		}
		// Иначе сохраняем как обычное слово
		state.tokens = append(state.tokens, Token{Type: WORD, Value: state.current.String()})
	}

	// Проверка незакрытых кавычек
	if state.inSingleQuote {
		return nil, fmt.Errorf("unclosed single quote")
	}
	if state.inDoubleQuote {
		return nil, fmt.Errorf("unclosed double quote")
	}

	return state.tokens, nil
}

// isValidVariableName проверяет, является ли строка корректным именем переменной.
// Имя переменной должно начинаться с буквы или подчеркивания и содержать только
// буквы, цифры и подчеркивания.
func (l *Lexer) isValidVariableName(name string) bool {
	// Пустая строка не является валидным именем переменной
	if len(name) == 0 {
		return false
	}

	runes := []rune(name)

	// Первый символ должен быть буквой или подчеркиванием
	if !unicode.IsLetter(runes[0]) && runes[0] != '_' {
		return false
	}

	// Остальные символы могут быть буквами, цифрами или подчеркиваниями
	for i := 1; i < len(runes); i++ {
		if !unicode.IsLetter(runes[i]) && !unicode.IsDigit(runes[i]) && runes[i] != '_' {
			return false
		}
	}

	return true
}
