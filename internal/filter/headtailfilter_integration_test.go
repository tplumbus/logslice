package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestHeadFilter_WithPipeline_LimitsOutput(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:01Z","msg":"first"}`,
		`{"ts":"2024-01-01T00:00:02Z","msg":"second"}`,
		`{"ts":"2024-01-01T00:00:03Z","msg":"third"}`,
		`{"ts":"2024-01-01T00:00:04Z","msg":"fourth"}`,
	}

	head, _ := NewHeadTailFilter("head", 2)

	var buf bytes.Buffer
	p := NewPipeline(strings.NewReader(strings.Join(lines, "\n")), &buf)
	p.AddFilter(head)
	if err := p.Run(); err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	got := strings.Split(out, "\n")
	if len(got) != 2 {
		t.Fatalf("expected 2 output lines, got %d: %v", len(got), got)
	}
	if !strings.Contains(got[0], "first") {
		t.Errorf("line 0: expected 'first', got %q", got[0])
	}
	if !strings.Contains(got[1], "second") {
		t.Errorf("line 1: expected 'second', got %q", got[1])
	}
}

func TestTailFilter_WithPipeline_FlushRequired(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:01Z","msg":"alpha"}`,
		`{"ts":"2024-01-01T00:00:02Z","msg":"beta"}`,
		`{"ts":"2024-01-01T00:00:03Z","msg":"gamma"}`,
		`{"ts":"2024-01-01T00:00:04Z","msg":"delta"}`,
		`{"ts":"2024-01-01T00:00:05Z","msg":"epsilon"}`,
	}

	tail, _ := NewHeadTailFilter("tail", 2)

	var buf bytes.Buffer
	p := NewPipeline(strings.NewReader(strings.Join(lines, "\n")), &buf)
	p.AddFilter(tail)
	if err := p.Run(); err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	// Pipeline won't emit tail lines; caller must flush manually.
	flushed := tail.Flush()
	if len(flushed) != 2 {
		t.Fatalf("expected 2 flushed lines, got %d", len(flushed))
	}
	if !strings.Contains(flushed[0], "delta") {
		t.Errorf("expected 'delta' in flushed[0], got %q", flushed[0])
	}
	if !strings.Contains(flushed[1], "epsilon") {
		t.Errorf("expected 'epsilon' in flushed[1], got %q", flushed[1])
	}
}
