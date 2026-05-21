package output

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_ToBuffer(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriterFromIO(&buf)

	lines := []string{`{"ts":"2024-01-01T00:00:00Z","msg":"hello"}`, `{"ts":"2024-01-02T00:00:00Z","msg":"world"}`}
	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	want := strings.Join(lines, "\n")
	if got != want {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", got, want)
	}
}

func TestWriter_ToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.log")
	w, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if err := w.WriteLine("test line"); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "test line") {
		t.Errorf("expected 'test line' in file, got: %q", string(data))
	}
}

func TestWriter_ToStdout(t *testing.T) {
	w, err := NewWriter("")
	if err != nil {
		t.Fatalf("unexpected error for stdout writer: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriter_InvalidPath(t *testing.T) {
	_, err := NewWriter("/no/such/dir/out.log")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
