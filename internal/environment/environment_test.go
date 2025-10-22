package environment

import (
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()

	if env.global == nil {
		t.Error("Global environment should be initialized")
	}
	if env.local == nil {
		t.Error("Local environment should be initialized")
	}

	// Проверяем, что системные переменные скопированы
	if len(env.global) == 0 {
		t.Error("Global environment should contain system variables")
	}
}

func TestEnvironment_SetAndGet(t *testing.T) {
	env := NewEnvironment()

	// Тест установки и получения локальной переменной
	env.Set("TEST_VAR", "test_value")
	value, exists := env.Get("TEST_VAR")

	if !exists {
		t.Error("Variable should exist after setting")
	}
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
}

func TestEnvironment_LocalOverridesGlobal(t *testing.T) {
	env := NewEnvironment()

	// Устанавливаем глобальную переменную
	env.global["OVERRIDE_TEST"] = "global_value"

	// Устанавливаем локальную с тем же именем
	env.Set("OVERRIDE_TEST", "local_value")

	value, exists := env.Get("OVERRIDE_TEST")
	if !exists {
		t.Error("Variable should exist")
	}
	if value != "local_value" {
		t.Errorf("Expected 'local_value', got '%s'", value)
	}
}

func TestEnvironment_GetNonExistent(t *testing.T) {
	env := NewEnvironment()

	_, exists := env.Get("NON_EXISTENT_VAR")
	if exists {
		t.Error("Non-existent variable should not exist")
	}
}

func TestEnvironment_Unset(t *testing.T) {
	env := NewEnvironment()

	env.Set("TO_UNSET", "value")
	env.Unset("TO_UNSET")

	_, exists := env.Get("TO_UNSET")
	if exists {
		t.Error("Variable should not exist after unset")
	}
}

func TestEnvironment_ClearLocal(t *testing.T) {
	env := NewEnvironment()

	env.Set("LOCAL1", "value1")
	env.Set("LOCAL2", "value2")
	env.ClearLocal()

	_, exists1 := env.Get("LOCAL1")
	_, exists2 := env.Get("LOCAL2")

	if exists1 || exists2 {
		t.Error("Local variables should be cleared")
	}
}

func TestEnvironment_GetAll(t *testing.T) {
	env := NewEnvironment()

	env.Set("LOCAL_VAR", "local_value")
	all := env.GetAll()

	// Проверяем, что есть как минимум локальная переменная
	found := false
	for _, envVar := range all {
		if envVar == "LOCAL_VAR=local_value" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Local variable should be in GetAll() result")
	}
}

func TestEnvironment_ListLocal(t *testing.T) {
	env := NewEnvironment()

	env.Set("LOCAL1", "value1")
	env.Set("LOCAL2", "value2")

	local := env.ListLocal()

	if len(local) != 2 {
		t.Errorf("Expected 2 local variables, got %d", len(local))
	}

	if local["LOCAL1"] != "value1" {
		t.Error("LOCAL1 should have correct value")
	}
	if local["LOCAL2"] != "value2" {
		t.Error("LOCAL2 should have correct value")
	}
}
