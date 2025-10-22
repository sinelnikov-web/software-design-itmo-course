package lexer

// TokenType определяет тип лексического токена.
// Используется для классификации различных элементов командной строки.
type TokenType int

const (
	WORD   TokenType = iota // Обычное слово или команда
	PIPE                    // Оператор пайплайна (|)
	SQUOTE                  // Одинарные кавычки (')
	DQUOTE                  // Двойные кавычки (")
	ASSIGN                  // Присваивание переменной (=)
)

// Token представляет лексический токен - минимальную единицу разбора.
// Содержит тип токена и его строковое значение.
type Token struct {
	Type  TokenType // Тип токена (WORD, PIPE, SQUOTE, DQUOTE, ASSIGN)
	Value string    // Строковое значение токена
}

// String возвращает строковое представление токена для отладки.
// Формат: "TYPE(value)" для токенов со значением, "TYPE" для токенов без значения.
func (t Token) String() string {
	switch t.Type {
	case WORD:
		return "WORD(" + t.Value + ")"
	case PIPE:
		return "PIPE"
	case SQUOTE:
		return "SQUOTE(" + t.Value + ")"
	case DQUOTE:
		return "DQUOTE(" + t.Value + ")"
	case ASSIGN:
		return "ASSIGN(" + t.Value + ")"
	default:
		return "UNKNOWN"
	}
}
