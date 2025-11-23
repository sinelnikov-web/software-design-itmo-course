package builtins

import (
	"testing"
)

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	// Проверяем, что все встроенные команды зарегистрированы
	expectedCommands := []string{"cat", "echo", "wc", "pwd", "exit"}

	for _, cmdName := range expectedCommands {
		if !registry.IsBuiltin(cmdName) {
			t.Errorf("Command %s should be registered", cmdName)
		}
	}
}

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

func TestRegistry_IsBuiltin(t *testing.T) {
	registry := NewRegistry()

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
			result := registry.IsBuiltin(tt.command)
			if result != tt.expected {
				t.Errorf("Registry.IsBuiltin(%s) = %v, expected %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()
	commands := registry.List()

	expectedCount := 5 // cat, echo, wc, pwd, exit
	if len(commands) != expectedCount {
		t.Errorf("Registry.List() returned %d commands, expected %d", len(commands), expectedCount)
	}

	// Проверяем, что все ожидаемые команды присутствуют
	expectedCommands := []string{"cat", "echo", "wc", "pwd", "exit"}
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
