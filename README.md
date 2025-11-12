# software-design-itmo-course
# GoCLI

[![Go Version](https://img.shields.io/badge/Go-1.25-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](https://opensource.org/licenses/MIT)

[![Build](https://github.com/sinelnikov-web/software-design-itmo-course/actions/workflows/build.yml/badge.svg)](https://github.com/sinelnikov-web/software-design-itmo-course/actions/workflows/build.yml)
[![Lint](https://github.com/sinelnikov-web/software-design-itmo-course/actions/workflows/lint.yml/badge.svg)](https://github.com/sinelnikov-web/software-design-itmo-course/actions/workflows/lint.yml)
[![Tests](https://github.com/sinelnikov-web/software-design-itmo-course/actions/workflows/test.yml/badge.svg)](https://github.com/sinelnikov-web/software-design-itmo-course/actions/workflows/test.yml)
Минималистичный интерпретатор командной строки, написанный на Go.

## Возможности

- **Встроенные команды**: `cat`, `echo`, `wc`, `pwd`, `exit`
- **Кавычки**: одинарные и двойные кавычки
- **Переменные окружения**: поддержка присваиваний `name=value`
- **Внешние программы**: вызов внешних программ, если команда не является встроенной
- **Потоки ввода-вывода**: поддержка stdin, stdout, stderr и кодов возврата

## Сборка и запуск

### Требования
- Go 1.25 или выше

### Сборка и запуск
```bash
# Клонирование репозитория
git clone https://github.com/username/gocli.git
cd gocli

# Сборка
make build

# Запуск
make run

# Или напрямую
go build -o gocli-cli ./cmd/gocli
./gocli-cli
```

### Тестирование
```bash
# Запуск тестов
make test

# Запуск линтера
make lint
```

## Примеры использования

```bash
# Встроенные команды
> echo "Hello, World!"
Hello, World!

> pwd
/home/user/project

> cat file.txt
Содержимое файла

> wc file.txt
1 3 18 file.txt

# Внешние команды
> ls -la
total 48
drwxr-xr-x 3 user user 4096 Oct 21 23:16 .
...

# Переменные окружения
> VAR=test echo hello
hello

# Кавычки
> echo "Hello World"
Hello World

> echo 'Hello World'
Hello World

# Команда exit
> exit          # Завершает shell с кодом 0
> exit 0        # Завершает shell с кодом 0
> exit 1        # Завершает shell с кодом 1
> exit 255      # Завершает shell с кодом 255
> exit "123"    # Завершает shell с кодом 123 (кавычки убираются)
> exit '123'    # Завершает shell с кодом 123 (кавычки убираются)
> exit abc      # Ошибка: "exit: abc: numeric argument required", код 2
> exit 256      # Код приводится к 0 (256 % 256)
> exit -1       # Код приводится к 255

# Примечание: при использовании в make или других скриптах
# ненулевой код возврата считается ошибкой:
# make run  # Если ввести exit 123, make получит код 123 и покажет ошибку
# Это нормальное поведение - shell корректно передает код возврата родительскому процессу
```

## Архитектура

Проект разделен на независимые компоненты:

1. **REPL** - основной цикл ввода-вывода
2. **Lexer** - токенизация командной строки
3. **Parser** - построение AST
4. **Executor** - выполнение команд
5. **Builtins Registry** - реестр встроенных команд
6. **Environment** - управление переменными окружения


## Структура проекта

```
gocli/
├── cmd/gocli/              # Точка входа
├── internal/
│   ├── shell/              # REPL логика
│   ├── lexer/              # Токенизация
│   ├── parser/             # Парсинг и AST
│   ├── executor/           # Выполнение команд
│   ├── builtins/           # Встроенные команды
│   └── environment/         # Управление переменными окружения
├── Makefile               # Команды сборки
└── README.md              # Документация
```

## Лицензия

MIT License