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
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("assignment without value")
			}

			nextToken := tokens[i+1]
			if nextToken.Type != lexer.WORD && nextToken.Type != lexer.SQUOTE && nextToken.Type != lexer.DQUOTE {
				return nil, fmt.Errorf("invalid assignment value")
			}

			assignment := &Assignment{
				Name:  token.Value,
				Value: p.createArgument(nextToken),
			}
			command.Assignments = append(command.Assignments, assignment)
			i++
		} else {
			if command.Name == "" {
				command.Name = token.Value
			} else {
				arg := p.createArgument(token)
				command.Args = append(command.Args, arg)
			}
		}
	}

	if command.Name == "" {
		return nil, fmt.Errorf("command name is required")
	}

	return command, nil
}

func (p *Parser) createArgument(token lexer.Token) *Argument {
	quoted := token.Type == lexer.SQUOTE || token.Type == lexer.DQUOTE
	return &Argument{
		Value:  token.Value,
		Quoted: quoted,
	}
}
