package builtins

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCdToExplicitDirectory(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	tmpDir := t.TempDir()

	cmd := NewCdCommand()
	var stderr bytes.Buffer

	exitCode := cmd.Execute(
		[]string{tmpDir},
		nil,
		nil,
		nil,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("cd returned non-zero exit code: %d, stderr=%s", exitCode, stderr.String())
	}

	newWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working directory: %v", err)
	}

	if newWd != tmpDir {
		t.Fatalf("expected wd=%s, got %s", tmpDir, newWd)
	}
}

func TestCdWithoutArgumentsGoesHome(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("cannot get home directory: %v", err)
	}

	cmd := NewCdCommand()
	var stderr bytes.Buffer

	exitCode := cmd.Execute(
		nil,
		nil,
		nil,
		nil,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("cd returned non-zero exit code: %d, stderr=%s", exitCode, stderr.String())
	}

	newWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working directory: %v", err)
	}

	if filepath.Clean(newWd) != filepath.Clean(home) {
		t.Fatalf("expected wd=%s, got %s", home, newWd)
	}
}

func TestCdTooManyArguments(t *testing.T) {
	cmd := NewCdCommand()
	var stderr bytes.Buffer

	exitCode := cmd.Execute(
		[]string{"a", "b"},
		nil,
		nil,
		nil,
		&stderr,
	)

	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code")
	}

	if stderr.Len() == 0 {
		t.Fatalf("expected error message in stderr")
	}
}
