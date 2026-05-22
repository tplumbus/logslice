package filter

import (
	"strings"
	"testing"
	"time"
)

func TestLimitFilter_WithPipeline_StopsEarly(t *testing.T) {
	now := time.Now().UTC()
	lines := []string{
		buildLine(now, "info", "first"),
		buildLine(now, "info", "second"),
		buildLine(now, "info", "third"),
		buildLine(now, "info", "fourth"),
		buildLine(now, "info", "fifth"),
	}

	lf, err := NewLimitFilter(3)
	if err != nil {
		t.Fatalf("NewLimitFilter: %v", err)
	}

	p := NewPipeline(lines, lf)
	results, err := p.Run()
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, want := range []string{"first", "second", "third"} {
		if !strings.Contains(results[i], want) {
			t.Errorf("result[%d] = %q, want to contain %q", i, results[i], want)
		}
	}
}

func TestLimitFilter_WithPipeline_LimitExceedsInput(t *testing.T) {
	now := time.Now().UTC()
	lines := []string{
		buildLine(now, "warn", "only"),
		buildLine(now, "warn", "two"),
	}

	lf, _ := NewLimitFilter(10)
	p := NewPipeline(lines, lf)
	results, err := p.Run()
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
