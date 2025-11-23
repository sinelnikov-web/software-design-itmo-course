package environment

import (
	"os"
	"strings"
)

const envVarParts = 2

// Environment управляет переменными окружения shell'а.
// Поддерживает как глобальные (системные), так и локальные переменные сессии.
type Environment struct {
	global map[string]string // Глобальные переменные (наследуются от системы)
	local  map[string]string // Локальные переменные сессии
}

// NewEnvironment создает новое окружение.
// Инициализирует глобальные переменные из системного окружения.
func NewEnvironment() *Environment {
	env := &Environment{
		global: make(map[string]string),
		local:  make(map[string]string),
	}

	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", envVarParts)
		if len(parts) == envVarParts {
			env.global[parts[0]] = parts[1]
		}
	}

	return env
}

// Set устанавливает переменную в локальном окружении.
// Локальные переменные имеют приоритет над глобальными.
func (e *Environment) Set(name, value string) {
	e.local[name] = value
}

// Get возвращает значение переменной.
// Сначала проверяет локальные, затем глобальные переменные.
func (e *Environment) Get(name string) (string, bool) {
	if value, exists := e.local[name]; exists {
		return value, true
	}
	if value, exists := e.global[name]; exists {
		return value, true
	}
	return "", false
}

// GetAll возвращает все переменные окружения для передачи внешним командам.
// Объединяет глобальные и локальные переменные (локальные имеют приоритет).
func (e *Environment) GetAll() []string {
	all := make(map[string]string)

	for k, v := range e.global {
		all[k] = v
	}

	for k, v := range e.local {
		all[k] = v
	}

	var result []string
	for k, v := range all {
		result = append(result, k+"="+v)
	}
	return result
}

// Unset удаляет переменную из локального окружения.
func (e *Environment) Unset(name string) {
	delete(e.local, name)
}

// ClearLocal очищает все локальные переменные.
func (e *Environment) ClearLocal() {
	e.local = make(map[string]string)
}

// ListLocal возвращает список локальных переменных.
func (e *Environment) ListLocal() map[string]string {
	result := make(map[string]string, len(e.local))
	for k, v := range e.local {
		result[k] = v
	}
	return result
}

// ListGlobal возвращает список глобальных переменных.
func (e *Environment) ListGlobal() map[string]string {
	result := make(map[string]string, len(e.global))
	for k, v := range e.global {
		result[k] = v
	}
	return result
}
