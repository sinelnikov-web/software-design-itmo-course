package parser

import (
	"testing"

	"gocli/internal/lexer"
)

// TestParser_Parse тестирует парсинг токенов в AST.
// Проверяет корректное построение AST для простых команд, команд с присваиваниями,
// команд с кавычками, пайплайнов и обработку ошибок для пустых команд.
func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []lexer.Token
		expected Node
		wantErr  bool
	}{
		{
			name: "simple command",
			tokens: []lexer.Token{
				{Type: lexer.WORD, Value: "echo"},
				{Type: lexer.WORD, Value: "hello"},
			},
			expected: &Command{
				Name: "echo",
				Args: []*Argument{{Value: "hello", Quoted: false}},
			},
			wantErr: false,
		},
		{
			name: "command with assignment",
			tokens: []lexer.Token{
				{Type: lexer.ASSIGN, Value: "VAR"},
				{Type: lexer.WORD, Value: "value"},
				{Type: lexer.WORD, Value: "echo"},
				{Type: lexer.WORD, Value: "hello"},
			},
			expected: &Command{
				Name:        "echo",
				Args:        []*Argument{{Value: "hello", Quoted: false}},
				Assignments: []*Assignment{{Name: "VAR", Value: &Argument{Value: "value", Quoted: false}}},
			},
			wantErr: false,
		},
		{
			name: "command with quotes",
			tokens: []lexer.Token{
				{Type: lexer.WORD, Value: "echo"},
				{Type: lexer.SQUOTE, Value: "hello world"},
			},
			expected: &Command{
				Name: "echo",
				Args: []*Argument{{Value: "hello world", Quoted: true}},
			},
			wantErr: false,
		},
		{
			name: "pipeline with two commands",
			tokens: []lexer.Token{
				{Type: lexer.WORD, Value: "echo"},
				{Type: lexer.WORD, Value: "hello"},
				{Type: lexer.PIPE, Value: "|"},
				{Type: lexer.WORD, Value: "wc"},
			},
			expected: &Pipeline{
				Commands: []*Command{
					{Name: "echo", Args: []*Argument{{Value: "hello", Quoted: false}}},
					{Name: "wc", Args: []*Argument{}},
				},
			},
			wantErr: false,
		},
		{
			name:     "empty command",
			tokens:   []lexer.Token{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "grep command",
			tokens: []lexer.Token{
				{Type: lexer.WORD, Value: "grep"},
				{Type: lexer.WORD, Value: "pattern"},
			},
			expected: &Command{
				Name: "grep",
				Args: []*Argument{{Value: "pattern", Quoted: false}},
			},
			wantErr: false,
		},
		{
			name: "grep with flag -w",
			tokens: []lexer.Token{
				{Type: lexer.WORD, Value: "grep"},
				{Type: lexer.WORD, Value: "-w"},
				{Type: lexer.WORD, Value: "pattern"},
			},
			expected: &Command{
				Name: "grep",
				Args: []*Argument{
					{Value: "-w", Quoted: false},
					{Value: "pattern", Quoted: false},
				},
			},
			wantErr: false,
		},
		{
			name: "grep with multiple flags",
			tokens: []lexer.Token{
				{Type: lexer.WORD, Value: "grep"},
				{Type: lexer.WORD, Value: "-i"},
				{Type: lexer.WORD, Value: "-w"},
				{Type: lexer.WORD, Value: "pattern"},
			},
			expected: &Command{
				Name: "grep",
				Args: []*Argument{
					{Value: "-i", Quoted: false},
					{Value: "-w", Quoted: false},
					{Value: "pattern", Quoted: false},
				},
			},
			wantErr: false,
		},
		{
			name: "pipeline with grep",
			tokens: []lexer.Token{
				{Type: lexer.WORD, Value: "echo"},
				{Type: lexer.WORD, Value: "hello"},
				{Type: lexer.PIPE, Value: "|"},
				{Type: lexer.WORD, Value: "grep"},
				{Type: lexer.WORD, Value: "hello"},
			},
			expected: &Pipeline{
				Commands: []*Command{
					{Name: "echo", Args: []*Argument{{Value: "hello", Quoted: false}}},
					{Name: "grep", Args: []*Argument{{Value: "hello", Quoted: false}}},
				},
			},
			wantErr: false,
		},
		{
			name: "pipeline with grep and flags",
			tokens: []lexer.Token{
				{Type: lexer.WORD, Value: "cat"},
				{Type: lexer.WORD, Value: "file"},
				{Type: lexer.PIPE, Value: "|"},
				{Type: lexer.WORD, Value: "grep"},
				{Type: lexer.WORD, Value: "-w"},
				{Type: lexer.WORD, Value: "pattern"},
			},
			expected: &Pipeline{
				Commands: []*Command{
					{Name: "cat", Args: []*Argument{{Value: "file", Quoted: false}}},
					{Name: "grep", Args: []*Argument{
						{Value: "-w", Quoted: false},
						{Value: "pattern", Quoted: false},
					}},
				},
			},
			wantErr: false,
		},
	}

	parser := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.tokens)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !compareNodes(result, tt.expected) {
					t.Errorf("Parser.Parse() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

func compareNodes(a, b Node) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch aCmd := a.(type) {
	case *Command:
		if bCmd, ok := b.(*Command); ok {
			return compareCommands(aCmd, bCmd)
		}
	case *Pipeline:
		if bPipeline, ok := b.(*Pipeline); ok {
			return comparePipelines(aCmd, bPipeline)
		}
	}

	return false
}

func compareCommands(a, b *Command) bool {
	if a.Name != b.Name {
		return false
	}
	if len(a.Args) != len(b.Args) {
		return false
	}
	for i, arg := range a.Args {
		if !compareArguments(arg, b.Args[i]) {
			return false
		}
	}
	if len(a.Assignments) != len(b.Assignments) {
		return false
	}
	for i, assignment := range a.Assignments {
		if !compareAssignments(assignment, b.Assignments[i]) {
			return false
		}
	}
	return true
}

func compareArguments(a, b *Argument) bool {
	return a.Value == b.Value && a.Quoted == b.Quoted
}

func compareAssignments(a, b *Assignment) bool {
	return a.Name == b.Name && compareArguments(a.Value, b.Value)
}

func comparePipelines(a, b *Pipeline) bool {
	if len(a.Commands) != len(b.Commands) {
		return false
	}
	for i, cmd := range a.Commands {
		if !compareCommands(cmd, b.Commands[i]) {
			return false
		}
	}
	return true
}
