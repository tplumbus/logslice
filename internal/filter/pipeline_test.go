package filter

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildLine(ts, level, msg string) string {
	return `{"time":"` + ts + `","level":"` + level + `","msg":"` + msg + `"}`
}

func TestPipeline_Run_TimeRangeOnly(t *testing.T) {
	tr := TimeRange{
		From: mustTime("2024-01-01T10:00:00Z"),
		To:   mustTime("2024-01-01T11:00:00Z"),
	}
	lines := []string{
		buildLine("2024-01-01T09:59:59Z", "info", "before"),
		buildLine("2024-01-01T10:30:00Z", "info", "inside"),
		buildLine("2024-01-01T11:00:01Z", "info", "after"),
	}
	input := strings.Join(lines, "\n")
	var buf bytes.Buffer
	p := NewPipeline(tr, nil)
	n, err := p.Run(strings.NewReader(input), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 line written, got %d", n)
	}
	if !strings.Contains(buf.String(), "inside") {
		t.Errorf("expected 'inside' in output, got: %s", buf.String())
	}
}

func TestPipeline_Run_FieldFilter(t *testing.T) {
	tr := TimeRange{From: time.Time{}} // open bounds
	lines := []string{
		buildLine("2024-01-01T10:00:00Z", "info", "one"),
		buildLine("2024-01-01T10:01:00Z", "error", "two"),
		buildLine("2024-01-01T10:02:00Z", "info", "three"),
	}
	input := strings.Join(lines, "\n")
	fq, _ := ParseFieldQuery("level=error")
	var buf bytes.Buffer
	p := NewPipeline(tr, []FieldQuery{fq})
	n, err := p.Run(strings.NewReader(input), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 line, got %d", n)
	}
	if !strings.Contains(buf.String(), "two") {
		t.Errorf("expected 'two' in output")
	}
}

func TestPipeline_Run_EmptyInput(t *testing.T) {
	p := NewPipeline(TimeRange{}, nil)
	var buf bytes.Buffer
	n, err := p.Run(strings.NewReader(""), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}

func TestPipeline_Run_SkipsInvalidJSON(t *testing.T) {
	p := NewPipeline(TimeRange{}, nil)
	var buf bytes.Buffer
	n, err := p.Run(strings.NewReader("not json\n"), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines for invalid JSON, got %d", n)
	}
}
