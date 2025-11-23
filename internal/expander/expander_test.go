package expander

import (
	"testing"

	"gocli/internal/environment"
	"gocli/internal/parser"
)

const (
	testValue   = "value"
	testVarName = "$VAR"
	testExit    = "exit"
	testTxtFile = "test.txt"
)

// TestExpander_ExpandVariable проверяет базовую подстановку переменных.
func TestExpander_ExpandVariable(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "$VAR",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != testValue {
		t.Errorf("expected '%s', got '%s'", testValue, expanded.Value)
	}
}

// TestExpander_ExpandVariableInDoubleQuotes проверяет подстановку переменных в двойных кавычках.
func TestExpander_ExpandVariableInDoubleQuotes(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "$VAR",
		Quoted:    true,
		QuoteType: parser.DoubleQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != testValue {
		t.Errorf("expected '%s', got '%s'", testValue, expanded.Value)
	}
}

// TestExpander_NoExpandInSingleQuotes проверяет, что подстановки не выполняются в одинарных кавычках.
func TestExpander_NoExpandInSingleQuotes(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "$VAR",
		Quoted:    true,
		QuoteType: parser.SingleQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != testVarName {
		t.Errorf("expected '%s', got '%s'", testVarName, expanded.Value)
	}
}

// TestExpander_ExpandMultipleVariables проверяет подстановку нескольких переменных.
func TestExpander_ExpandMultipleVariables(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("X", "ex")
	env.Set("Y", "it")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "$X$Y",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != testExit {
		t.Errorf("expected '%s', got '%s'", testExit, expanded.Value)
	}
}

// TestExpander_ExpandVariableWithText проверяет подстановку переменной в тексте.
func TestExpander_ExpandVariableWithText(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "prefix_$VAR_suffix",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != "prefix_value_suffix" {
		t.Errorf("expected 'prefix_value_suffix', got '%s'", expanded.Value)
	}
}

// TestExpander_ExpandUndefinedVariable проверяет подстановку несуществующей переменной.
func TestExpander_ExpandUndefinedVariable(t *testing.T) {
	env := environment.NewEnvironment()
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "$UNDEFINED",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != "" {
		t.Errorf("expected empty string, got '%s'", expanded.Value)
	}
}

// TestExpander_ExpandVariableInBraces проверяет подстановку переменной в фигурных скобках.
func TestExpander_ExpandVariableInBraces(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "${VAR}",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != testValue {
		t.Errorf("expected '%s', got '%s'", testValue, expanded.Value)
	}
}

// TestExpander_ExpandCommand проверяет подстановку переменных в команде.
func TestExpander_ExpandCommand(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("X", "ex")
	env.Set("Y", "it")
	exp := NewExpander(env)

	cmd := &parser.Command{
		Name: "$X$Y",
		Args: []*parser.Argument{
			{
				Value:     "$X",
				Quoted:    false,
				QuoteType: parser.NoQuote,
			},
		},
		Assignments: []*parser.Assignment{},
	}

	expanded, err := exp.Expand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expandedCmd, ok := expanded.(*parser.Command)
	if !ok {
		t.Fatalf("expected *parser.Command, got %T", expanded)
	}

	if expandedCmd.Name != testExit {
		t.Errorf("expected command name '%s', got '%s'", testExit, expandedCmd.Name)
	}

	if len(expandedCmd.Args) != 1 || expandedCmd.Args[0].Value != "ex" {
		t.Errorf("expected arg 'ex', got %v", expandedCmd.Args)
	}
}

// TestExpander_ExpandPipeline проверяет подстановку переменных в пайплайне.
func TestExpander_ExpandPipeline(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("FILE", "test.txt")
	exp := NewExpander(env)

	pipeline := &parser.Pipeline{
		Commands: []*parser.Command{
			{
				Name: "cat",
				Args: []*parser.Argument{
					{
						Value:     "$FILE",
						Quoted:    false,
						QuoteType: parser.NoQuote,
					},
				},
			},
			{
				Name: "wc",
				Args: []*parser.Argument{},
			},
		},
	}

	expanded, err := exp.Expand(pipeline)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expandedPipeline, ok := expanded.(*parser.Pipeline)
	if !ok {
		t.Fatalf("expected *parser.Pipeline, got %T", expanded)
	}

	if len(expandedPipeline.Commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(expandedPipeline.Commands))
	}

	if expandedPipeline.Commands[0].Args[0].Value != testTxtFile {
		t.Errorf("expected arg '%s', got '%s'", testTxtFile, expandedPipeline.Commands[0].Args[0].Value)
	}
}

// TestExpander_ExpandAssignment проверяет подстановку переменных в присваиваниях.
func TestExpander_ExpandAssignment(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("BASE", "test")
	exp := NewExpander(env)

	cmd := &parser.Command{
		Name: "echo",
		Args: []*parser.Argument{},
		Assignments: []*parser.Assignment{
			{
				Name: "FILE",
				Value: &parser.Argument{
					Value:     "$BASE.txt",
					Quoted:    false,
					QuoteType: parser.NoQuote,
				},
			},
		},
	}

	expanded, err := exp.Expand(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expandedCmd, ok := expanded.(*parser.Command)
	if !ok {
		t.Fatalf("expected *parser.Command, got %T", expanded)
	}

	if len(expandedCmd.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(expandedCmd.Assignments))
	}

	if expandedCmd.Assignments[0].Value.Value != testTxtFile {
		t.Errorf("expected assignment value '%s', got '%s'", testTxtFile, expandedCmd.Assignments[0].Value.Value)
	}
}

// TestExpander_ExpandDollarSign проверяет обработку символа $ не как переменной.
func TestExpander_ExpandDollarSign(t *testing.T) {
	env := environment.NewEnvironment()
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "$",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != "$" {
		t.Errorf("expected '$', got '%s'", expanded.Value)
	}
}

// TestExpander_ExpandDollarAtEnd проверяет обработку $ в конце строки.
func TestExpander_ExpandDollarAtEnd(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "text$",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded.Value != "text$" {
		t.Errorf("expected 'text$', got '%s'", expanded.Value)
	}
}

// TestExpander_ExpandEscapedDollar проверяет обработку экранированного $.
func TestExpander_ExpandEscapedDollar(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "\\$VAR",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Экранированный $ должен остаться как $, переменная не подставляется
	if expanded.Value != testVarName {
		t.Errorf("expected '%s', got '%s'", testVarName, expanded.Value)
	}
}

// TestExpander_ExpandDoubleBackslash проверяет обработку двойного обратного слеша перед $.
func TestExpander_ExpandDoubleBackslash(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "\\\\$VAR",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// \\$VAR: первый \ экранирует второй \, поэтому $ не экранирован
	// Результат: \ + значение VAR = \value
	if expanded.Value != "\\value" {
		t.Errorf("expected '\\value', got '%s'", expanded.Value)
	}
}

// TestExpander_ExpandTripleBackslash проверяет обработку тройного обратного слеша перед $.
func TestExpander_ExpandTripleBackslash(t *testing.T) {
	env := environment.NewEnvironment()
	env.Set("VAR", "value")
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "\\\\\\$VAR",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// \\\$VAR: первый \ экранирует второй \, третий \ экранирует $
	// Результат: \\ + $VAR = \\$VAR
	if expanded.Value != "\\\\$VAR" {
		t.Errorf("expected '\\\\$VAR', got '%s'", expanded.Value)
	}
}

// TestExpander_ExpandInvalidVariableName проверяет обработку некорректного имени переменной.
func TestExpander_ExpandInvalidVariableName(t *testing.T) {
	env := environment.NewEnvironment()
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "$123",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	expanded, err := exp.expandArgument(arg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// $123 не является валидным именем переменной, поэтому $ остается как есть
	if expanded.Value != "$123" {
		t.Errorf("expected '$123', got '%s'", expanded.Value)
	}
}

// TestExpander_ExpandUnclosedBraces проверяет обработку незакрытых фигурных скобок.
func TestExpander_ExpandUnclosedBraces(t *testing.T) {
	env := environment.NewEnvironment()
	exp := NewExpander(env)

	arg := &parser.Argument{
		Value:     "${VAR",
		Quoted:    false,
		QuoteType: parser.NoQuote,
	}

	_, err := exp.expandArgument(arg)
	if err == nil {
		t.Error("expected error for unclosed ${ variable")
	}
}
