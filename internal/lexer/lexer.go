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
	var tokens []Token
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	runes := []rune(input)

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		switch {
		case char == '\'' && !inDoubleQuote:
			if inSingleQuote {
				tokens = append(tokens, Token{Type: SQUOTE, Value: current.String()})
				current.Reset()
				inSingleQuote = false
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, Token{Type: WORD, Value: current.String()})
					current.Reset()
				}
				inSingleQuote = true
			}

		case char == '"' && !inSingleQuote:
			if inDoubleQuote {
				tokens = append(tokens, Token{Type: DQUOTE, Value: current.String()})
				current.Reset()
				inDoubleQuote = false
			} else {
				if current.Len() > 0 {
					tokens = append(tokens, Token{Type: WORD, Value: current.String()})
					current.Reset()
				}
				inDoubleQuote = true
			}

		case char == '|' && !inSingleQuote && !inDoubleQuote:
			if current.Len() > 0 {
				tokens = append(tokens, Token{Type: WORD, Value: current.String()})
				current.Reset()
			}
			tokens = append(tokens, Token{Type: PIPE, Value: "|"})

		case unicode.IsSpace(char) && !inSingleQuote && !inDoubleQuote:
			if current.Len() > 0 {
				tokens = append(tokens, Token{Type: WORD, Value: current.String()})
				current.Reset()
			}

		case char == '=' && !inSingleQuote && !inDoubleQuote:
			if current.Len() > 0 {
				word := current.String()
				if l.isValidVariableName(word) {
					tokens = append(tokens, Token{Type: ASSIGN, Value: word})
					current.Reset()
				} else {
					current.WriteRune(char)
				}
			} else {
				current.WriteRune(char)
			}

		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		if inSingleQuote || inDoubleQuote {
			return nil, fmt.Errorf("unclosed quote")
		}
		tokens = append(tokens, Token{Type: WORD, Value: current.String()})
	}

	if inSingleQuote {
		return nil, fmt.Errorf("unclosed single quote")
	}
	if inDoubleQuote {
		return nil, fmt.Errorf("unclosed double quote")
	}

	return tokens, nil
}

// isValidVariableName проверяет, является ли строка корректным именем переменной.
// Имя переменной должно начинаться с буквы или подчеркивания и содержать только
// буквы, цифры и подчеркивания.
func (l *Lexer) isValidVariableName(name string) bool {
	if len(name) == 0 {
		return false
	}

	runes := []rune(name)

	if !unicode.IsLetter(runes[0]) && runes[0] != '_' {
		return false
	}

	for i := 1; i < len(runes); i++ {
		if !unicode.IsLetter(runes[i]) && !unicode.IsDigit(runes[i]) && runes[i] != '_' {
			return false
		}
	}

	return true
}
