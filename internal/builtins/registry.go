package builtins

import (
	"fmt"
)

// Registry управляет реестром встроенных команд.
// Позволяет регистрировать, получать и проверять наличие команд.
type Registry struct {
	commands map[string]Builtin // Карта команд: имя -> реализация
}

// NewRegistry создает новый реестр встроенных команд.
// Автоматически регистрирует все доступные встроенные команды.
func NewRegistry() *Registry {
	registry := &Registry{
		commands: make(map[string]Builtin),
	}

	registry.Register(NewCatCommand())
	registry.Register(NewEchoCommand())
	registry.Register(NewWcCommand())
	registry.Register(NewPwdCommand())
	registry.Register(NewExitCommand())

	return registry
}

// Register регистрирует новую встроенную команду в реестре.
// Если команда с таким именем уже существует, она будет перезаписана.
func (r *Registry) Register(command Builtin) {
	r.commands[command.Name()] = command
}

// Get возвращает встроенную команду по имени.
// Возвращает команду и флаг существования.
func (r *Registry) Get(name string) (Builtin, bool) {
	command, exists := r.commands[name]
	return command, exists
}

// List возвращает список имен всех зарегистрированных команд.
func (r *Registry) List() []string {
	var names []string
	for name := range r.commands {
		names = append(names, name)
	}
	return names
}

// IsBuiltin проверяет, является ли команда встроенной.
// Возвращает true, если команда зарегистрирована в реестре.
func (r *Registry) IsBuiltin(name string) bool {
	_, exists := r.commands[name]
	return exists
}

// String возвращает строковое представление реестра для отладки.
// Формат: "Registry with N commands: [cmd1, cmd2, ...]"
func (r *Registry) String() string {
	return fmt.Sprintf("Registry with %d commands: %v", len(r.commands), r.List())
}
