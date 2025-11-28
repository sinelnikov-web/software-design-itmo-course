package builtins

import (
	"bytes"
	"strings"
	"testing"
)

func TestGrepCommand_Name(t *testing.T) {
	command := NewGrepCommand()
	expected := "grep"

	if command.Name() != expected {
		t.Errorf("GrepCommand.Name() = %s, expected %s", command.Name(), expected)
	}
}

func TestGrepCommand_Execute_BasicSearch(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("line one\nline two\nline three\n")

	exitCode := command.Execute([]string{"two"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "two") {
		t.Errorf("GrepCommand.Execute() output = %q, expected to contain 'two'", output)
	}
}

func TestGrepCommand_Execute_RegexSearch(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("line 1\nline 2\nline 10\n")

	// Поиск по регулярному выражению: строки, заканчивающиеся на 0
	exitCode := command.Execute([]string{"0$"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "10") {
		t.Errorf("GrepCommand.Execute() output = %q, expected to contain '10'", output)
	}
	if strings.Contains(output, "1\n") || strings.Contains(output, "2\n") {
		t.Errorf("GrepCommand.Execute() output = %q, should not contain '1' or '2'", output)
	}
}

func TestGrepCommand_Execute_CaseInsensitive(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("Line One\nLINE TWO\nline three\n")

	exitCode := command.Execute([]string{"-i", "two"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "TWO") {
		t.Errorf("GrepCommand.Execute() output = %q, expected to contain 'TWO'", output)
	}
}

func TestGrepCommand_Execute_WordBoundary(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("test\ntesting\ntest case\n")

	// Поиск слова "test" целиком
	exitCode := command.Execute([]string{"-w", "test"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	// Должны найтись "test" и "test case", но не "testing"
	if !strings.Contains(output, "test\n") && !strings.Contains(output, "test case") {
		t.Errorf("GrepCommand.Execute() output = %q, expected to contain 'test' or 'test case'", output)
	}
	if strings.Contains(output, "testing") {
		t.Errorf("GrepCommand.Execute() output = %q, should not contain 'testing'", output)
	}
}

func TestGrepCommand_Execute_AfterContext(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("line 1\nline 2\nmatch\nline 4\nline 5\nline 6\n")

	exitCode := command.Execute([]string{"-A", "2", "match"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Должны быть: match, line 4, line 5
	if len(lines) < 3 {
		t.Errorf("GrepCommand.Execute() expected at least 3 lines, got %d", len(lines))
	}

	hasMatch := false
	hasLine4 := false
	hasLine5 := false

	for _, line := range lines {
		if strings.Contains(line, "match") {
			hasMatch = true
		}
		if strings.Contains(line, "line 4") {
			hasLine4 = true
		}
		if strings.Contains(line, "line 5") {
			hasLine5 = true
		}
	}

	if !hasMatch {
		t.Errorf("GrepCommand.Execute() output should contain 'match'")
	}
	if !hasLine4 {
		t.Errorf("GrepCommand.Execute() output should contain 'line 4'")
	}
	if !hasLine5 {
		t.Errorf("GrepCommand.Execute() output should contain 'line 5'")
	}
}

func TestGrepCommand_Execute_AfterContextZero(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("line 1\nmatch\nline 3\n")

	exitCode := command.Execute([]string{"-A", "0", "match"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	// Должна быть только строка с match
	if !strings.Contains(output, "match") {
		t.Errorf("GrepCommand.Execute() output = %q, expected to contain 'match'", output)
	}
	if strings.Contains(output, "line 3") {
		t.Errorf("GrepCommand.Execute() output = %q, should not contain 'line 3'", output)
	}
}

func TestGrepCommand_Execute_OverlappingContext(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Два совпадения близко друг к другу
	stdin.WriteString("line 1\nmatch1\nline 3\nmatch2\nline 5\n")

	exitCode := command.Execute([]string{"-A", "2", "match"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	// Оба совпадения должны быть найдены
	if !strings.Contains(output, "match1") {
		t.Errorf("GrepCommand.Execute() output should contain 'match1'")
	}
	if !strings.Contains(output, "match2") {
		t.Errorf("GrepCommand.Execute() output should contain 'match2'")
	}
}

func TestGrepCommand_Execute_NoMatch(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("line one\nline two\nline three\n")

	exitCode := command.Execute([]string{"nonexistent"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 1 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 1 (no match)", exitCode)
	}

	output := stdout.String()
	if output != "" {
		t.Errorf("GrepCommand.Execute() output = %q, expected empty", output)
	}
}

func TestGrepCommand_Execute_InvalidRegex(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("test\n")

	exitCode := command.Execute([]string{"[invalid"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 2 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 2 (error)", exitCode)
	}

	errorOutput := stderr.String()
	if !strings.Contains(errorOutput, "invalid regular expression") {
		t.Errorf("GrepCommand.Execute() stderr = %q, expected to contain 'invalid regular expression'", errorOutput)
	}
}

func TestGrepCommand_Execute_NoPattern(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := command.Execute([]string{}, nil, &stdin, &stdout, &stderr)

	if exitCode != 2 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 2 (error)", exitCode)
	}

	errorOutput := stderr.String()
	if !strings.Contains(errorOutput, "pattern required") {
		t.Errorf("GrepCommand.Execute() stderr = %q, expected to contain 'pattern required'", errorOutput)
	}
}

func TestGrepCommand_Execute_CombinedFlags(t *testing.T) {
	command := NewGrepCommand()

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stdin.WriteString("Line One\nLINE TWO\nline three\nTEST\n")

	// Комбинация -i и -w
	exitCode := command.Execute([]string{"-i", "-w", "test"}, nil, &stdin, &stdout, &stderr)

	if exitCode != 0 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 0", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "TEST") {
		t.Errorf("GrepCommand.Execute() output = %q, expected to contain 'TEST'", output)
	}
}

func TestGrepCommand_Execute_WithFile(t *testing.T) {
	// Создаем временный файл для теста
	// В реальном тесте можно использовать ioutil.TempFile, но для простоты используем существующий подход
	command := NewGrepCommand()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Используем несуществующий файл для проверки обработки ошибок
	exitCode := command.Execute([]string{"pattern", "nonexistent.txt"}, nil, nil, &stdout, &stderr)

	if exitCode != 1 {
		t.Errorf("GrepCommand.Execute() exitCode = %d, expected 1 (file error)", exitCode)
	}

	errorOutput := stderr.String()
	if !strings.Contains(errorOutput, "nonexistent.txt") {
		t.Errorf("GrepCommand.Execute() stderr = %q, expected to contain 'nonexistent.txt'", errorOutput)
	}
}
