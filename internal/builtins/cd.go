package builtins

import (
	"fmt"
	"io"
	"os"
)

type CdCommand struct{}

func NewCdCommand() *CdCommand { return &CdCommand{} }

func (c *CdCommand) Name() string { return "cd" }

func (c *CdCommand) Execute(args []string, _ map[string]string, _ io.Reader, _ io.Writer, stderr io.Writer) int {
	var target string
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(stderr, "cd: cannot determine home directory: %v\n", err)
			return 1
		}
		target = home
	} else if len(args) == 1 {
		target = args[0]
	} else {
		fmt.Fprintln(stderr, "cd: too many arguments")
		return 1
	}

	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(stderr, "cd: %v\n", err)
		return 1
	}
	return 0
}
