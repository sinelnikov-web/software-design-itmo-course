package builtins

import (
	"bytes"
	"strings"
	"testing"
)

// TestWcCommand_Execute тестирует выполнение команды wc с различными входными данными.
// Проверяет подсчет строк, слов и байт для различных сценариев: пустой ввод,
// одна строка, несколько строк с пробелами, и обработку ошибок для несуществующих файлов.
func TestWcCommand_Execute(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "count from stdin",
			args:     []string{},
			input:    "hello world\nline 2\nline 3",
			expected: "3 6 26",
			wantErr:  false,
		},
		{
			name:     "empty input",
			args:     []string{},
			input:    "",
			expected: "0 0 0",
			wantErr:  false,
		},
		{
			name:     "single line",
			args:     []string{},
			input:    "hello",
			expected: "1 1 6",
			wantErr:  false,
		},
		{
			name:     "multiple lines with spaces",
			args:     []string{},
			input:    "hello world\nthis is a test\nfinal line",
			expected: "3 8 38",
			wantErr:  false,
		},
		{
			name:     "file that doesn't exist",
			args:     []string{"nonexistent.txt"},
			input:    "",
			expected: "",
			wantErr:  true,
		},
	}

	command := NewWcCommand()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			var stdin bytes.Buffer

			stdin.WriteString(tt.input)

			exitCode := command.Execute(tt.args, nil, &stdin, &stdout, &stderr)

			if tt.wantErr {
				if exitCode == 0 {
					t.Errorf("WcCommand.Execute() expected error, got exit code 0")
				}
			} else {
				if exitCode != 0 {
					t.Errorf("WcCommand.Execute() exitCode = %d, expected 0", exitCode)
				}

				output := strings.TrimSpace(stdout.String())
				if !strings.Contains(output, tt.expected) {
					t.Errorf("WcCommand.Execute() output = %q, expected to contain %q", output, tt.expected)
				}
			}
		})
	}
}

// TestWcCommand_Count тестирует внутренний метод count команды wc.
// Проверяет корректность подсчета строк, слов и байт для различных входных данных.
func TestWcCommand_Count(t *testing.T) {
	command := NewWcCommand()

	tests := []struct {
		name          string
		input         string
		expectedLines int
		expectedWords int
		expectedBytes int
	}{
		{
			name:          "empty input",
			input:         "",
			expectedLines: 0,
			expectedWords: 0,
			expectedBytes: 0,
		},
		{
			name:          "single line",
			input:         "hello world",
			expectedLines: 1,
			expectedWords: 2,
			expectedBytes: 12,
		},
		{
			name:          "multiple lines",
			input:         "line 1\nline 2\nline 3",
			expectedLines: 3,
			expectedWords: 6,
			expectedBytes: 21,
		},
		{
			name:          "lines with different word counts",
			input:         "hello\nworld test\nfinal",
			expectedLines: 3,
			expectedWords: 4,
			expectedBytes: 23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			lines, words, bytes := command.count(reader)

			if lines != tt.expectedLines {
				t.Errorf("WcCommand.count() lines = %d, expected %d", lines, tt.expectedLines)
			}
			if words != tt.expectedWords {
				t.Errorf("WcCommand.count() words = %d, expected %d", words, tt.expectedWords)
			}
			if bytes != tt.expectedBytes {
				t.Errorf("WcCommand.count() bytes = %d, expected %d", bytes, tt.expectedBytes)
			}
		})
	}
}

// TestWcCommand_Name тестирует получение имени команды wc.
// Проверяет, что команда возвращает корректное имя "wc".
func TestWcCommand_Name(t *testing.T) {
	command := NewWcCommand()
	expected := "wc"

	if command.Name() != expected {
		t.Errorf("WcCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}
