package filter

import (
	"strings"
	"testing"
)

func TestCompareFilter_WithPipeline_FiltersByThreshold(t *testing.T) {
	lines := []string{
		`{"level":"info","latency":50}`,
		`{"level":"warn","latency":150}`,
		`{"level":"error","latency":300}`,
		`{"level":"debug","latency":10}`,
	}

	f, err := NewCompareFilter("latency", OpGte, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := NewPipeline(strings.NewReader(strings.Join(lines, "\n")), f)
	var results []string
	err = p.Run(func(line string) error {
		results = append(results, line)
		return nil
	})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d: %v", len(results), results)
	}
	if !strings.Contains(results[0], "150") {
		t.Errorf("expected first result to contain latency 150, got %s", results[0])
	}
	if !strings.Contains(results[1], "300") {
		t.Errorf("expected second result to contain latency 300, got %s", results[1])
	}
}

func TestCompareFilter_WithPipeline_ParsedFromFlag(t *testing.T) {
	lines := []string{
		`{"code":200,"path":"/ok"}`,
		`{"code":404,"path":"/missing"}`,
		`{"code":500,"path":"/error"}`,
	}

	f, err := ParseCompareFlag("code:gte:400")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := NewPipeline(strings.NewReader(strings.Join(lines, "\n")), f)
	var results []string
	err = p.Run(func(line string) error {
		results = append(results, line)
		return nil
	})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
