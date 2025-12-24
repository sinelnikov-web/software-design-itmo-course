package builtins

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLsExplicitDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(tmpDir+"/file1.txt", []byte("hello"), 0644)
	if err != nil {
		t.Fatalf("cannot create file: %v", err)
	}

	err = os.Mkdir(tmpDir+"/dir1", 0755)
	if err != nil {
		t.Fatalf("cannot create dir: %v", err)
	}

	cmd := NewLsCommand()
	var stdout, stderr bytes.Buffer

	exitCode := cmd.Execute(
		[]string{tmpDir},
		nil,
		nil,
		&stdout,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("ls returned non-zero exit code: %d, stderr=%s", exitCode, stderr.String())
	}

	output := stdout.String()

	if !strings.Contains(output, "file1.txt") {
		t.Fatalf("expected file1.txt in output, got:\n%s", output)
	}

	if !strings.Contains(output, "dir1") {
		t.Fatalf("expected dir1 in output, got:\n%s", output)
	}
}

func TestLsCurrentDirectory(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get wd: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	tmpDir := t.TempDir()
	err = os.WriteFile(tmpDir+"/a.txt", []byte("a"), 0644)
	if err != nil {
		t.Fatalf("cannot create file: %v", err)
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("cannot chdir: %v", err)
	}

	cmd := NewLsCommand()
	var stdout, stderr bytes.Buffer

	exitCode := cmd.Execute(
		nil,
		nil,
		nil,
		&stdout,
		&stderr,
	)

	if exitCode != 0 {
		t.Fatalf("ls returned non-zero exit code: %d, stderr=%s", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), "a.txt") {
		t.Fatalf("expected a.txt in output, got:\n%s", stdout.String())
	}
}

func TestLsTooManyArguments(t *testing.T) {
	cmd := NewLsCommand()
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
