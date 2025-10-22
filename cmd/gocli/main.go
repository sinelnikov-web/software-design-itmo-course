package main

import (
	"fmt"
	"os"

	"gocli/internal/shell"
)

func main() {
	sh := shell.NewShell()

	if err := sh.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
