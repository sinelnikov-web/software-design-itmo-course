package builtins

import (
	"testing"
)

// TestRegistry_Register тестирует регистрацию встроенных команд в реестре.
// Проверяет, что все стандартные команды (cat, echo, wc, pwd, exit, grep) зарегистрированы.
func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	// Проверяем, что все встроенные команды зарегистрированы
	expectedCommands := []string{"cat", "echo", "wc", "pwd", "exit", "grep"}

	for _, cmdName := range expectedCommands {
		if !registry.IsBuiltin(cmdName) {
			t.Errorf("Command %s should be registered", cmdName)
		}
	}
}

// TestRegistry_Get тестирует получение команды из реестра по имени.
// Проверяет корректность работы для существующих и несуществующих команд.
func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()

	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"existing command", "echo", true},
		{"existing command", "cat", true},
		{"non-existing command", "nonexistent", false},
		{"empty command", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := registry.Get(tt.command)
			if exists != tt.expected {
				t.Errorf("Registry.Get(%s) exists = %v, expected %v", tt.command, exists, tt.expected)
			}
		})
	}
}

// TestRegistry_IsBuiltin тестирует проверку, является ли команда встроенной.
// Проверяет корректность определения встроенных и внешних команд.
func TestRegistry_IsBuiltin(t *testing.T) {
	registry := NewRegistry()

	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"builtin command", "echo", true},
		{"builtin command", "cat", true},
		{"non-builtin command", "ping", false},
		{"empty command", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := registry.IsBuiltin(tt.command)
			if result != tt.expected {
				t.Errorf("Registry.IsBuiltin(%s) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

// TestRegistry_List тестирует получение списка всех зарегистрированных команд.
// Проверяет, что список содержит все ожидаемые встроенные команды.
func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()
	commands := registry.List()

	expectedCount := 8 // cat, echo, wc, pwd, exit, grep
	if len(commands) != expectedCount {
		t.Errorf("Registry.List() returned %d commands, expected %d", len(commands), expectedCount)
	}

	// Проверяем, что все ожидаемые команды присутствуют
	expectedCommands := []string{"cat", "echo", "wc", "pwd", "exit", "grep"}
	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range commands {
			if cmd == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Command %s not found in registry list", expected)
		}
	}
}
