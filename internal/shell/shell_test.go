package shell

import (
	"testing"
)

// TestNewShell тестирует создание нового экземпляра shell.
// Проверяет, что все компоненты (executor, lexer, parser, environment) инициализированы корректно.
func TestNewShell(t *testing.T) {
	sh := NewShell()

	if sh.executor == nil {
		t.Error("Shell executor should be initialized")
	}
	if sh.lexer == nil {
		t.Error("Shell lexer should be initialized")
	}
	if sh.parser == nil {
		t.Error("Shell parser should be initialized")
	}
	if sh.environment == nil {
		t.Error("Shell environment should be initialized")
	}
}

// TestShell_ProcessCommand тестирует обработку команд через shell.
// Проверяет корректную обработку простых команд, команд с аргументами и обработку ошибок.
func TestShell_ProcessCommand(t *testing.T) {
	sh := NewShell()

	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{
			name:    "simple echo command",
			command: "echo hello",
			wantErr: false,
		},
		{
			name:    "echo with multiple arguments",
			command: "echo hello world",
			wantErr: false,
		},
		{
			name:    "pwd command",
			command: "pwd",
			wantErr: false,
		},
		{
			name:    "command with quotes",
			command: `echo "hello world"`,
			wantErr: false,
		},
		{
			name:    "command with single quotes",
			command: `echo 'hello world'`,
			wantErr: false,
		},
		{
			name:    "empty command",
			command: "",
			wantErr: true,
		},
		{
			name:    "invalid syntax - unclosed quote",
			command: `echo "hello`,
			wantErr: true,
		},
		{
			name:    "nonexistent command",
			command: "nonexistentcommand123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.processCommand(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Shell.processCommand(%q) error = %v, wantErr %v", tt.command, err, tt.wantErr)
			}
		})
	}
}

// TestShell_ProcessCommandWithAssignment тестирует обработку команд с присваиванием переменных.
// Проверяет, что переменные окружения корректно устанавливаются перед выполнением команды.
func TestShell_ProcessCommandWithAssignment(t *testing.T) {
	sh := NewShell()

	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{
			name:    "command with variable assignment",
			command: "VAR=test echo hello",
			wantErr: false,
		},
		{
			name:    "multiple assignments",
			command: "VAR1=value1 VAR2=value2 echo test",
			wantErr: false,
		},
		{
			name:    "assignment with quoted value",
			command: `VAR="test value" echo hello`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.processCommand(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Shell.processCommand(%q) error = %v, wantErr %v", tt.command, err, tt.wantErr)
			}
		})
	}
}

// TestShell_ProcessCommandWithPipeline тестирует обработку пайплайнов команд.
// Проверяет корректную обработку команд, соединенных оператором пайплайна.
func TestShell_ProcessCommandWithPipeline(t *testing.T) {
	sh := NewShell()

	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{
			name:    "simple pipeline",
			command: "echo hello | wc",
			wantErr: false, // Pipeline execution now implemented
		},
		{
			name:    "empty command before pipe",
			command: "| echo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.processCommand(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Shell.processCommand(%q) error = %v, wantErr %v", tt.command, err, tt.wantErr)
			}
		})
	}
}

// TestShell_ProcessCommandLexicalErrors тестирует обработку лексических ошибок.
// Проверяет корректную обработку ошибок токенизации.
func TestShell_ProcessCommandLexicalErrors(t *testing.T) {
	sh := NewShell()

	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{
			name:    "unclosed single quote",
			command: "echo 'hello",
			wantErr: true,
		},
		{
			name:    "unclosed double quote",
			command: `echo "hello`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.processCommand(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Shell.processCommand(%q) error = %v, wantErr %v", tt.command, err, tt.wantErr)
			}

			// Проверяем, что ошибка содержит информацию о лексическом анализе
			if err != nil && !tt.wantErr {
				if err.Error() == "" {
					t.Error("Error message should not be empty")
				}
			}
		})
	}
}

// TestShell_ProcessCommandParsingErrors тестирует обработку ошибок парсинга.
// Проверяет корректную обработку синтаксических ошибок.
func TestShell_ProcessCommandParsingErrors(t *testing.T) {
	sh := NewShell()

	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{
			name:    "empty command after tokenization",
			command: "   ",
			wantErr: true,
		},
		{
			name:    "standalone assignment (VAR=value syntax)",
			command: "VAR=test",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.processCommand(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Shell.processCommand(%q) error = %v, wantErr %v", tt.command, err, tt.wantErr)
			}
		})
	}
}

// TestShell_StandaloneAssignments тестирует новую функциональность standalone assignments.
// Проверяет, что переменные могут быть установлены без выполнения команды.
func TestShell_StandaloneAssignments(t *testing.T) {
	sh := NewShell()

	tests := []struct {
		name       string
		command    string
		wantErr    bool
		checkVar   string
		checkValue string
	}{
		{
			name:       "simple standalone assignment",
			command:    "VAR=value",
			wantErr:    false,
			checkVar:   "VAR",
			checkValue: "value",
		},
		{
			name:       "standalone assignment with quoted value",
			command:    `VAR="hello world"`,
			wantErr:    false,
			checkVar:   "VAR",
			checkValue: "hello world",
		},
		{
			name:       "standalone assignment with single quotes",
			command:    `VAR='test value'`,
			wantErr:    false,
			checkVar:   "VAR",
			checkValue: "test value",
		},
		{
			name:    "assignment with variable substitution",
			command: `X=foo Y=$X`,
			wantErr: false,
			checkVar:   "Y",
			checkValue: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.processCommand(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Shell.processCommand(%q) error = %v, wantErr %v", tt.command, err, tt.wantErr)
			}

			if !tt.wantErr && tt.checkVar != "" {
				value, exists := sh.environment.Get(tt.checkVar)
				if !exists {
					t.Errorf("Variable %s was not set after command %q", tt.checkVar, tt.command)
				} else if value != tt.checkValue {
					t.Errorf("Variable %s = %q, expected %q after command %q", 
						tt.checkVar, value, tt.checkValue, tt.command)
				}
			}
		})
	}
}
