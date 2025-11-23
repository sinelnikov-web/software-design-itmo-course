package environment

import (
	"os"
	"strings"
	"sync"
)

const envVarParts = 2

// Environment управляет переменными окружения shell'а.
// Поддерживает как глобальные (системные), так и локальные переменные сессии.
type Environment struct {
	mu     sync.RWMutex      // Мьютекс для защиты доступа к maps
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
func (env *Environment) Set(name, value string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	env.local[name] = value
}

// Get возвращает значение переменной.
// Сначала проверяет локальные, затем глобальные переменные.
func (env *Environment) Get(name string) (string, bool) {
	env.mu.RLock()
	defer env.mu.RUnlock()
	if value, exists := env.local[name]; exists {
		return value, true
	}
	if value, exists := env.global[name]; exists {
		return value, true
	}
	return "", false
}

// GetAll возвращает все переменные окружения для передачи внешним командам.
// Объединяет глобальные и локальные переменные (локальные имеют приоритет).
func (env *Environment) GetAll() []string {
	env.mu.RLock()
	defer env.mu.RUnlock()

	all := make(map[string]string)

	for k, v := range env.global {
		all[k] = v
	}

	for k, v := range env.local {
		all[k] = v
	}

	var result []string
	for k, v := range all {
		result = append(result, k+"="+v)
	}
	return result
}

// Unset удаляет переменную из локального окружения.
func (env *Environment) Unset(name string) {
	env.mu.Lock()
	defer env.mu.Unlock()
	delete(env.local, name)
}

// ClearLocal очищает все локальные переменные.
func (env *Environment) ClearLocal() {
	env.mu.Lock()
	defer env.mu.Unlock()
	env.local = make(map[string]string)
}

// ListLocal возвращает список локальных переменных.
func (env *Environment) ListLocal() map[string]string {
	env.mu.RLock()
	defer env.mu.RUnlock()
	result := make(map[string]string, len(env.local))
	for k, v := range env.local {
		result[k] = v
	}
	return result
}

// ListGlobal возвращает список глобальных переменных.
func (env *Environment) ListGlobal() map[string]string {
	env.mu.RLock()
	defer env.mu.RUnlock()
	result := make(map[string]string, len(env.global))
	for k, v := range env.global {
		result[k] = v
	}
	return result
}

// GetAllMap возвращает все переменные окружения в виде map[string]string.
// Объединяет глобальные и локальные переменные (локальные имеют приоритет).
// Используется для передачи переменных окружения во встроенные команды.
func (env *Environment) GetAllMap() map[string]string {
	env.mu.RLock()
	defer env.mu.RUnlock()

	result := make(map[string]string)

	// Сначала копируем глобальные переменные
	for k, v := range env.global {
		result[k] = v
	}

	// Затем перезаписываем локальными (локальные имеют приоритет)
	for k, v := range env.local {
		result[k] = v
	}

	return result
}

// HasLocal проверяет, существует ли переменная в локальном окружении.
func (env *Environment) HasLocal(name string) bool {
	env.mu.RLock()
	defer env.mu.RUnlock()
	_, exists := env.local[name]
	return exists
}

// HasGlobal проверяет, существует ли переменная в глобальном окружении.
func (env *Environment) HasGlobal(name string) bool {
	env.mu.RLock()
	defer env.mu.RUnlock()
	_, exists := env.global[name]
	return exists
}

// GetLocal возвращает значение переменной из локального окружения.
// Возвращает значение и флаг существования.
func (env *Environment) GetLocal(name string) (string, bool) {
	env.mu.RLock()
	defer env.mu.RUnlock()
	value, exists := env.local[name]
	return value, exists
}
