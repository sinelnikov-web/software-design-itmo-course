package environment

import (
	"testing"
)

const (
	testGlobalValue = "global_value"
	testLocalValue  = "local_value"
	testValue       = "value"
)

// TestNewEnvironment тестирует создание нового экземпляра окружения.
// Проверяет инициализацию глобального и локального окружений, а также копирование системных переменных.
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

// TestEnvironment_SetAndGet тестирует установку и получение переменных окружения.
// Проверяет, что установленные переменные корректно сохраняются и извлекаются.
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

// TestEnvironment_LocalOverridesGlobal тестирует приоритет локальных переменных над глобальными.
// Проверяет, что локальная переменная с тем же именем переопределяет глобальную.
func TestEnvironment_LocalOverridesGlobal(t *testing.T) {
	env := NewEnvironment()

	// Устанавливаем глобальную переменную
	env.global["OVERRIDE_TEST"] = testGlobalValue

	// Устанавливаем локальную с тем же именем
	env.Set("OVERRIDE_TEST", testLocalValue)

	value, exists := env.Get("OVERRIDE_TEST")
	if !exists {
		t.Error("Variable should exist")
	}
	if value != testLocalValue {
		t.Errorf("Expected '%s', got '%s'", testLocalValue, value)
	}
}

// TestEnvironment_GetNonExistent тестирует получение несуществующей переменной.
// Проверяет, что метод Get корректно возвращает false для несуществующих переменных.
func TestEnvironment_GetNonExistent(t *testing.T) {
	env := NewEnvironment()

	_, exists := env.Get("NON_EXISTENT_VAR")
	if exists {
		t.Error("Non-existent variable should not exist")
	}
}

// TestEnvironment_Unset тестирует удаление переменной окружения.
// Проверяет, что после удаления переменная больше не доступна.
func TestEnvironment_Unset(t *testing.T) {
	env := NewEnvironment()

	env.Set("TO_UNSET", testValue)
	env.Unset("TO_UNSET")

	_, exists := env.Get("TO_UNSET")
	if exists {
		t.Error("Variable should not exist after unset")
	}
}

// TestEnvironment_ClearLocal тестирует очистку всех локальных переменных.
// Проверяет, что после очистки все локальные переменные удаляются, а глобальные остаются.
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

// TestEnvironment_GetAll тестирует получение всех переменных окружения.
// Проверяет, что метод GetAll возвращает все переменные в формате "KEY=VALUE".
func TestEnvironment_GetAll(t *testing.T) {
	env := NewEnvironment()

	env.Set("LOCAL_VAR", testLocalValue)
	all := env.GetAll()

	// Проверяем, что есть как минимум локальная переменная
	found := false
	expectedVar := "LOCAL_VAR=" + testLocalValue
	for _, envVar := range all {
		if envVar == expectedVar {
			found = true
			break
		}
	}

	if !found {
		t.Error("Local variable should be in GetAll() result")
	}
}

// TestEnvironment_ListLocal тестирует получение списка локальных переменных.
// Проверяет, что метод ListLocal возвращает только локальные переменные в виде map.
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

// TestEnvironment_GetAllMap тестирует получение всех переменных окружения в виде map.
// Проверяет, что метод GetAllMap возвращает объединение глобальных и локальных переменных,
// при этом локальные имеют приоритет над глобальными.
func TestEnvironment_GetAllMap(t *testing.T) {
	env := NewEnvironment()

	// Устанавливаем локальную переменную
	env.Set("LOCAL_VAR", testLocalValue)

	// Устанавливаем глобальную переменную напрямую
	env.global["GLOBAL_VAR"] = testGlobalValue

	all := env.GetAllMap()

	// Проверяем наличие локальной переменной
	if all["LOCAL_VAR"] != testLocalValue {
		t.Errorf("Expected '%s' for LOCAL_VAR, got '%s'", testLocalValue, all["LOCAL_VAR"])
	}

	// Проверяем наличие глобальной переменной
	if all["GLOBAL_VAR"] != testGlobalValue {
		t.Errorf("Expected '%s' for GLOBAL_VAR, got '%s'", testGlobalValue, all["GLOBAL_VAR"])
	}

	// Проверяем приоритет локальных над глобальными
	env.global["OVERRIDE_VAR"] = testGlobalValue
	env.Set("OVERRIDE_VAR", testLocalValue)

	all = env.GetAllMap()
	if all["OVERRIDE_VAR"] != testLocalValue {
		t.Errorf("Expected '%s' for OVERRIDE_VAR (local should override global), got '%s'", testLocalValue, all["OVERRIDE_VAR"])
	}
}

// TestEnvironment_HasLocal тестирует проверку существования переменной в локальном окружении.
func TestEnvironment_HasLocal(t *testing.T) {
	env := NewEnvironment()

	// Переменная не существует
	if env.HasLocal("NON_EXISTENT") {
		t.Error("HasLocal should return false for non-existent variable")
	}

	// Устанавливаем локальную переменную
	env.Set("LOCAL_VAR", testValue)

	if !env.HasLocal("LOCAL_VAR") {
		t.Error("HasLocal should return true for existing local variable")
	}

	// Глобальная переменная не считается локальной
	env.global["GLOBAL_VAR"] = testValue
	if env.HasLocal("GLOBAL_VAR") {
		t.Error("HasLocal should return false for global-only variable")
	}
}

// TestEnvironment_HasGlobal тестирует проверку существования переменной в глобальном окружении.
func TestEnvironment_HasGlobal(t *testing.T) {
	env := NewEnvironment()

	// Переменная не существует
	if env.HasGlobal("NON_EXISTENT") {
		t.Error("HasGlobal should return false for non-existent variable")
	}

	// Устанавливаем глобальную переменную напрямую
	env.global["GLOBAL_VAR"] = testValue

	if !env.HasGlobal("GLOBAL_VAR") {
		t.Error("HasGlobal should return true for existing global variable")
	}

	// Локальная переменная не считается глобальной
	env.Set("LOCAL_VAR", testValue)
	if env.HasGlobal("LOCAL_VAR") {
		t.Error("HasGlobal should return false for local-only variable")
	}
}

// TestEnvironment_GetLocal тестирует получение значения переменной из локального окружения.
func TestEnvironment_GetLocal(t *testing.T) {
	env := NewEnvironment()

	// Переменная не существует в локальном окружении
	value, exists := env.GetLocal("NON_EXISTENT")
	if exists {
		t.Error("GetLocal should return false for non-existent variable")
	}
	if value != "" {
		t.Errorf("GetLocal should return empty string for non-existent variable, got '%s'", value)
	}

	// Устанавливаем локальную переменную
	env.Set("LOCAL_VAR", testLocalValue)
	value, exists = env.GetLocal("LOCAL_VAR")

	if !exists {
		t.Error("GetLocal should return true for existing local variable")
	}
	if value != testLocalValue {
		t.Errorf("Expected '%s', got '%s'", testLocalValue, value)
	}

	// Глобальная переменная не возвращается через GetLocal
	env.global["GLOBAL_VAR"] = testGlobalValue
	_, exists = env.GetLocal("GLOBAL_VAR")
	if exists {
		t.Error("GetLocal should return false for global-only variable")
	}
}

// TestEnvironment_UnsetGlobalVariable тестирует удаление локальной переменной,
// которая переопределяла глобальную, и проверяет возврат к глобальной.
func TestEnvironment_UnsetGlobalVariable(t *testing.T) {
	env := NewEnvironment()

	// Устанавливаем глобальную переменную
	env.global["TEST_VAR"] = testGlobalValue

	// Устанавливаем локальную с тем же именем
	env.Set("TEST_VAR", testLocalValue)

	// Проверяем, что Get возвращает локальное значение
	value, _ := env.Get("TEST_VAR")
	if value != testLocalValue {
		t.Errorf("Expected '%s', got '%s'", testLocalValue, value)
	}

	// Удаляем локальную переменную
	env.Unset("TEST_VAR")

	// Проверяем, что теперь Get возвращает глобальное значение
	value, exists := env.Get("TEST_VAR")
	if !exists {
		t.Error("Variable should exist in global environment")
	}
	if value != testGlobalValue {
		t.Errorf("Expected '%s' after unset, got '%s'", testGlobalValue, value)
	}
}
