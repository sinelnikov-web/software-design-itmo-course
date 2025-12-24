package builtins

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type LsCommand struct{}

func NewLsCommand() *LsCommand { return &LsCommand{} }

func (l *LsCommand) Name() string { return "ls" }

func (l *LsCommand) Execute(args []string, _ map[string]string, _ io.Reader, stdout io.Writer, stderr io.Writer) int {
	var target string
	if len(args) == 0 {
		target = "."
	} else if len(args) == 1 {
		target = args[0]
	} else {
		fmt.Fprintln(stderr, "ls: too many arguments")
		return 1
	}

	info, err := os.Stat(target)
	if err != nil {
		fmt.Fprintf(stderr, "ls: %v\n", err)
		return 2
	}
	if !info.IsDir() {
		fmt.Fprintln(stdout, filepath.Base(target))
		return 0
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		fmt.Fprintf(stderr, "ls: %v\n", err)
		return 2
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintln(stdout, n)
	}
	return 0
}
