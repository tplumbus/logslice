package input

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestNewLineReader_SingleFile(t *testing.T) {
	path := writeTempFile(t, "{\"ts\":\"2024-01-01T00:00:00Z\"}\n{\"ts\":\"2024-01-02T00:00:00Z\"}\n")
	lr, err := NewLineReader([]string{path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var lines []string
	for l := range lr.Lines() {
		lines = append(lines, l)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestNewLineReader_MultipleFiles(t *testing.T) {
	p1 := writeTempFile(t, "line1\nline2\n")
	p2 := writeTempFile(t, "line3\n")
	lr, err := NewLineReader([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var lines []string
	for l := range lr.Lines() {
		lines = append(lines, l)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestNewLineReader_MissingFile(t *testing.T) {
	_, err := NewLineReader([]string{filepath.Join(t.TempDir(), "no-such-file.log")})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNewLineReader_SkipsBlankLines(t *testing.T) {
	path := writeTempFile(t, "line1\n\nline2\n\n")
	lr, err := NewLineReader([]string{path})
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for l := range lr.Lines() {
		lines = append(lines, l)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 non-blank lines, got %d", len(lines))
	}
}

func TestNewLineReader_Stdin(t *testing.T) {
	lr, err := NewLineReader(nil)
	if err != nil {
		t.Fatalf("unexpected error for stdin reader: %v", err)
	}
	if lr == nil {
		t.Fatal("expected non-nil reader")
	}
}
