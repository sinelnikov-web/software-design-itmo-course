package lexer

import (
	"testing"
)

func TestLexer_Tokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
		wantErr  bool
	}{
		{
			name:  "simple command",
			input: "echo hello",
			expected: []Token{
				{Type: WORD, Value: "echo"},
				{Type: WORD, Value: "hello"},
			},
			wantErr: false,
		},
		{
			name:  "command with single quotes",
			input: "echo 'hello world'",
			expected: []Token{
				{Type: WORD, Value: "echo"},
				{Type: SQUOTE, Value: "hello world"},
			},
			wantErr: false,
		},
		{
			name:  "command with double quotes",
			input: `echo "hello world"`,
			expected: []Token{
				{Type: WORD, Value: "echo"},
				{Type: DQUOTE, Value: "hello world"},
			},
			wantErr: false,
		},
		{
			name:  "command with pipe",
			input: "echo hello | wc",
			expected: []Token{
				{Type: WORD, Value: "echo"},
				{Type: WORD, Value: "hello"},
				{Type: PIPE, Value: "|"},
				{Type: WORD, Value: "wc"},
			},
			wantErr: false,
		},
		{
			name:  "assignment",
			input: "VAR=value echo hello",
			expected: []Token{
				{Type: ASSIGN, Value: "VAR"},
				{Type: WORD, Value: "value"},
				{Type: WORD, Value: "echo"},
				{Type: WORD, Value: "hello"},
			},
			wantErr: false,
		},
		{
			name:     "unclosed single quote",
			input:    "echo 'hello",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "unclosed double quote",
			input:    `echo "hello`,
			expected: nil,
			wantErr:  true,
		},
	}

	lexer := NewLexer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.Tokenize(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Lexer.Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(tokens) != len(tt.expected) {
					t.Errorf("Lexer.Tokenize() returned %d tokens, expected %d", len(tokens), len(tt.expected))
					return
				}

				for i, token := range tokens {
					if token.Type != tt.expected[i].Type {
						t.Errorf("Token %d: type = %v, expected %v", i, token.Type, tt.expected[i].Type)
					}
					if token.Value != tt.expected[i].Value {
						t.Errorf("Token %d: value = %v, expected %v", i, token.Value, tt.expected[i].Value)
					}
				}
			}
		})
	}
}

func TestLexer_isValidVariableName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid variable", "VAR", true},
		{"valid variable with underscore", "VAR_NAME", true},
		{"valid variable with numbers", "VAR123", true},
		{"invalid variable starting with number", "123VAR", false},
		{"invalid variable with special chars", "VAR-NAME", false},
		{"empty variable", "", false},
	}

	lexer := NewLexer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lexer.isValidVariableName(tt.input)
			if result != tt.expected {
				t.Errorf("isValidVariableName(%s) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
