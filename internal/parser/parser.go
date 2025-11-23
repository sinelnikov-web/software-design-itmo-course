package parser

import (
	"fmt"

	"gocli/internal/lexer"
)

// Parser выполняет синтаксический анализ токенов.
// Строит абстрактное синтаксическое дерево (AST) из последовательности токенов.
type Parser struct{}

// NewParser создает новый экземпляр парсера.
// Возвращает готовую к использованию структуру Parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse выполняет синтаксический анализ последовательности токенов.
// Строит AST, представляющий команды и пайплайны.
// Возвращает корневой узел AST и ошибку при некорректном синтаксисе.
func (p *Parser) Parse(tokens []lexer.Token) (Node, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	commands, err := p.parsePipeline(tokens)
	if err != nil {
		return nil, err
	}

	if len(commands) == 1 {
		return commands[0], nil
	}

	return &Pipeline{Commands: commands}, nil
}

func (p *Parser) parsePipeline(tokens []lexer.Token) ([]*Command, error) {
	var commands []*Command
	var currentTokens []lexer.Token

	for _, token := range tokens {
		if token.Type == lexer.PIPE {
			if len(currentTokens) == 0 {
				return nil, fmt.Errorf("empty command before pipe")
			}

			command, err := p.parseCommand(currentTokens)
			if err != nil {
				return nil, err
			}
			commands = append(commands, command)
			currentTokens = nil
		} else {
			currentTokens = append(currentTokens, token)
		}
	}

	if len(currentTokens) > 0 {
		command, err := p.parseCommand(currentTokens)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}

	return commands, nil
}

func (p *Parser) parseCommand(tokens []lexer.Token) (*Command, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	command := &Command{
		Args:        []*Argument{},
		Assignments: []*Assignment{},
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if token.Type == lexer.ASSIGN {
			assignment, skip, err := p.parseAssignment(tokens, i)
			if err != nil {
				return nil, err
			}
			command.Assignments = append(command.Assignments, assignment)
			i += skip
		} else {
			if err := p.addTokenToCommand(command, token); err != nil {
				return nil, err
			}
		}
	}

	// Команда может состоять только из assignments (например, x=5)
	// или должна иметь имя команды (например, x=5 echo hello)
	if command.Name == "" && len(command.Assignments) == 0 {
		return nil, fmt.Errorf("command name or assignment is required")
	}

	return command, nil
}

// parseAssignment обрабатывает присваивание переменной из токенов.
// Возвращает созданное присваивание, количество пропущенных токенов и ошибку.
func (p *Parser) parseAssignment(tokens []lexer.Token, i int) (*Assignment, int, error) {
	if i+1 >= len(tokens) {
		return nil, 0, fmt.Errorf("assignment without value")
	}

	nextToken := tokens[i+1]
	if nextToken.Type != lexer.WORD && nextToken.Type != lexer.SQUOTE && nextToken.Type != lexer.DQUOTE {
		return nil, 0, fmt.Errorf("invalid assignment value")
	}

	assignment := &Assignment{
		Name:  tokens[i].Value,
		Value: p.createArgument(nextToken),
	}
	return assignment, 1, nil
}

// addTokenToCommand добавляет токен к команде (как имя команды или аргумент).
func (p *Parser) addTokenToCommand(command *Command, token lexer.Token) error {
	if command.Name == "" {
		command.Name = token.Value
	} else {
		arg := p.createArgument(token)
		command.Args = append(command.Args, arg)
	}
	return nil
}

func (p *Parser) createArgument(token lexer.Token) *Argument {
	var quoteType QuoteType
	quoted := false

	switch token.Type {
	case lexer.SQUOTE:
		quoted = true
		quoteType = SingleQuote
	case lexer.DQUOTE:
		quoted = true
		quoteType = DoubleQuote
	default:
		quoteType = NoQuote
	}

	return &Argument{
		Value:     token.Value,
		Quoted:    quoted,
		QuoteType: quoteType,
	}
}
