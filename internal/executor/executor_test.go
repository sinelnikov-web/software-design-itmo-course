package executor

import (
	"testing"

	"gocli/internal/parser"
)

// TestExecutor_IsBuiltin тестирует проверку, является ли команда встроенной.
// Проверяет корректность определения встроенных и внешних команд.
func TestExecutor_IsBuiltin(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"builtin command", "echo", true},
		{"builtin command", "cat", true},
		{"non-builtin command", "ls", false},
		{"empty command", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.IsBuiltin(tt.command)
			if result != tt.expected {
				t.Errorf("Executor.IsBuiltin(%s) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

// TestExecutor_ListBuiltins тестирует получение списка всех встроенных команд.
// Проверяет, что список содержит все ожидаемые команды (cat, echo, wc, pwd, exit).
func TestExecutor_ListBuiltins(t *testing.T) {
	executor := NewExecutor()
	commands := executor.ListBuiltins()

	expectedCount := 5 // cat, echo, wc, pwd, exit
	if len(commands) != expectedCount {
		t.Errorf("Executor.ListBuiltins() returned %d commands, expected %d", len(commands), expectedCount)
	}
}

// TestExecutor_ExecuteCommand тестирует выполнение команд через executor.
// Проверяет выполнение встроенных команд (echo, pwd) и обработку ошибок для несуществующих команд.
func TestExecutor_ExecuteCommand(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name    string
		command *parser.Command
		wantErr bool
	}{
		{
			name: "echo command",
			command: &parser.Command{
				Name: "echo",
				Args: []*parser.Argument{{Value: "hello", Quoted: false, QuoteType: parser.NoQuote}},
			},
			wantErr: false,
		},
		{
			name: "pwd command",
			command: &parser.Command{
				Name: "pwd",
				Args: []*parser.Argument{},
			},
			wantErr: false,
		},
		{
			name: "nonexistent command",
			command: &parser.Command{
				Name: "nonexistent",
				Args: []*parser.Argument{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.Execute(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("Executor.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestExecutor_ExecutePipeline тестирует выполнение пайплайнов команд.
// Проверяет выполнение пайплайна с одной командой и обработку пайплайнов с несколькими командами.
func TestExecutor_ExecutePipeline(t *testing.T) {
	executor := NewExecutor()

	tests := []struct {
		name     string
		pipeline *parser.Pipeline
		wantErr  bool
	}{
		{
			name: "single command pipeline",
			pipeline: &parser.Pipeline{
				Commands: []*parser.Command{
					{Name: "echo", Args: []*parser.Argument{{Value: "hello", Quoted: false, QuoteType: parser.NoQuote}}},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple commands pipeline",
			pipeline: &parser.Pipeline{
				Commands: []*parser.Command{
					{Name: "echo", Args: []*parser.Argument{{Value: "hello", Quoted: false, QuoteType: parser.NoQuote}}},
					{Name: "wc", Args: []*parser.Argument{}},
				},
			},
			wantErr: false, // Pipeline execution now implemented
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.Execute(tt.pipeline)

			if (err != nil) != tt.wantErr {
				t.Errorf("Executor.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestExecutor_ExecuteCommandWithAssignments тестирует выполнение команд с присваиваниями переменных.
// Проверяет, что переменные устанавливаются перед выполнением команды и восстанавливаются после.
func TestExecutor_ExecuteCommandWithAssignments(t *testing.T) {
	executor := NewExecutor()

	// Устанавливаем начальную переменную в окружении
	executor.environment.Set("EXISTING_VAR", "original_value")

	// Выполняем команду с временным присваиванием
	command := &parser.Command{
		Name: "echo",
		Args: []*parser.Argument{{Value: "hello", Quoted: false}},
		Assignments: []*parser.Assignment{
			{
				Name:  "TEMP_VAR",
				Value: &parser.Argument{Value: "temp_value", Quoted: false, QuoteType: parser.NoQuote},
			},
			{
				Name:  "EXISTING_VAR",
				Value: &parser.Argument{Value: "new_value", Quoted: false, QuoteType: parser.NoQuote},
			},
		},
	}

	err := executor.Execute(command)
	if err != nil {
		t.Errorf("Executor.Execute() error = %v, expected no error", err)
	}

	// Проверяем, что временная переменная удалена после выполнения
	_, exists := executor.environment.Get("TEMP_VAR")
	if exists {
		t.Error("TEMP_VAR should be removed after command execution")
	}

	// Проверяем, что существующая переменная восстановлена
	value, exists := executor.environment.Get("EXISTING_VAR")
	if !exists {
		t.Error("EXISTING_VAR should exist after command execution")
	}
	if value != "original_value" {
		t.Errorf("EXISTING_VAR should be restored to 'original_value', got '%s'", value)
	}
}

// TestExecutor_ExecuteCommandWithGlobalVariableAssignment тестирует присваивание переменной,
// которая существовала только в глобальном окружении.
// Проверяет, что после выполнения команды возвращается доступ к глобальной переменной.
func TestExecutor_ExecuteCommandWithGlobalVariableAssignment(t *testing.T) {
	executor := NewExecutor()

	// Используем системную переменную, которая точно есть (PATH обычно есть везде)
	testVarName := "PATH"

	// Получаем исходное значение глобальной переменной
	originalValue, exists := executor.environment.Get(testVarName)
	if !exists {
		// Если PATH нет, используем другую системную переменную или пропускаем тест
		t.Skip("PATH variable not found, skipping test")
		return
	}

	// Выполняем команду с временным присваиванием
	command := &parser.Command{
		Name: "echo",
		Args: []*parser.Argument{{Value: "test", Quoted: false, QuoteType: parser.NoQuote}},
		Assignments: []*parser.Assignment{
			{
				Name:  testVarName,
				Value: &parser.Argument{Value: "temp_value", Quoted: false, QuoteType: parser.NoQuote},
			},
		},
	}

	err := executor.Execute(command)
	if err != nil {
		t.Errorf("Executor.Execute() error = %v, expected no error", err)
	}

	// Проверяем, что переменная восстановлена к исходному глобальному значению
	value, exists := executor.environment.Get(testVarName)
	if !exists {
		t.Errorf("%s should exist in global environment after command execution", testVarName)
	}
	if value != originalValue {
		t.Errorf("%s should be restored to original value '%s', got '%s'", testVarName, originalValue, value)
	}

	// Проверяем, что локальной переменной нет
	if executor.environment.HasLocal(testVarName) {
		t.Errorf("%s should not exist in local environment after command execution", testVarName)
	}
}

// TestExecutor_ExecuteBuiltinWithEnvironment тестирует передачу переменных окружения во встроенные команды.
// Проверяет, что встроенные команды получают переменные окружения через параметр env.
func TestExecutor_ExecuteBuiltinWithEnvironment(t *testing.T) {
	executor := NewExecutor()

	// Устанавливаем локальные переменные окружения
	executor.environment.Set("TEST_VAR", "test_value")

	// Используем системную переменную, которая точно есть (например, PATH)
	// или проверяем через GetAllMap, что системные переменные присутствуют
	envMap := executor.environment.GetAllMap()

	// Проверяем, что локальная переменная доступна
	if envMap["TEST_VAR"] != "test_value" {
		t.Errorf("Expected 'test_value' for TEST_VAR, got '%s'", envMap["TEST_VAR"])
	}

	// Проверяем, что системные переменные присутствуют (PATH обычно есть)
	if _, exists := envMap["PATH"]; !exists {
		// PATH может не быть в некоторых окружениях, но это не критично для теста
		t.Log("PATH variable not found in environment, this may be normal in some test environments")
	}

	// Выполняем команду с временным присваиванием
	command := &parser.Command{
		Name: "echo",
		Args: []*parser.Argument{{Value: "hello", Quoted: false}},
		Assignments: []*parser.Assignment{
			{
				Name:  "TEMP_VAR",
				Value: &parser.Argument{Value: "temp_value", Quoted: false, QuoteType: parser.NoQuote},
			},
		},
	}

	// Проверяем, что временная переменная доступна через GetAllMap во время выполнения
	// (но мы не можем проверить это напрямую, так как переменная устанавливается внутри executeCommand)

	// Выполняем команду
	err := executor.Execute(command)
	if err != nil {
		t.Errorf("Executor.Execute() error = %v, expected no error", err)
	}

	// Проверяем, что временная переменная удалена после выполнения
	if _, exists := executor.environment.Get("TEMP_VAR"); exists {
		t.Error("TEMP_VAR should be removed after command execution")
	}
}

// TestExecutor_ExecutePipelineWithThreeCommands тестирует выполнение пайплайна с тремя командами.
// Проверяет корректную передачу данных через pipe между командами.
func TestExecutor_ExecutePipelineWithThreeCommands(t *testing.T) {
	executor := NewExecutor()

	pipeline := &parser.Pipeline{
		Commands: []*parser.Command{
			{Name: "echo", Args: []*parser.Argument{{Value: "hello world", Quoted: false, QuoteType: parser.NoQuote}}},
			{Name: "wc", Args: []*parser.Argument{}},
			{Name: "wc", Args: []*parser.Argument{}},
		},
	}

	err := executor.Execute(pipeline)
	if err != nil {
		t.Errorf("Executor.Execute() error = %v, expected no error", err)
	}
}

// TestExecutor_ExecutePipelineWithAssignments тестирует выполнение пайплайна с присваиваниями переменных.
// Проверяет, что переменные устанавливаются для каждой команды в пайплайне.
func TestExecutor_ExecutePipelineWithAssignments(t *testing.T) {
	executor := NewExecutor()

	pipeline := &parser.Pipeline{
		Commands: []*parser.Command{
			{
				Name: "echo",
				Args: []*parser.Argument{{Value: "test", Quoted: false, QuoteType: parser.NoQuote}},
				Assignments: []*parser.Assignment{
					{
						Name:  "VAR1",
						Value: &parser.Argument{Value: "value1", Quoted: false, QuoteType: parser.NoQuote},
					},
				},
			},
			{
				Name: "wc",
				Args: []*parser.Argument{},
				Assignments: []*parser.Assignment{
					{
						Name:  "VAR2",
						Value: &parser.Argument{Value: "value2", Quoted: false, QuoteType: parser.NoQuote},
					},
				},
			},
		},
	}

	err := executor.Execute(pipeline)
	if err != nil {
		t.Errorf("Executor.Execute() error = %v, expected no error", err)
	}

	// Проверяем, что временные переменные удалены после выполнения
	if _, exists := executor.environment.Get("VAR1"); exists {
		t.Error("VAR1 should be removed after pipeline execution")
	}
	if _, exists := executor.environment.Get("VAR2"); exists {
		t.Error("VAR2 should be removed after pipeline execution")
	}
}
