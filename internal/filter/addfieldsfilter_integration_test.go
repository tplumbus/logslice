package filter

import (
	"strings"
	"testing"
)

func TestAddFieldsFilter_WithPipeline_AddsFields(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"started"}`,
		`{"level":"warn","msg":"slow query"}`,
	}

	af, err := NewAddFieldsFilter([]string{"env=prod", "app=logslice"})
	if err != nil {
		t.Fatalf("NewAddFieldsFilter: %v", err)
	}

	p := NewPipeline(lines)
	p.AddTransformer(af)

	results, err := p.Run()
	if err != nil {
		t.Fatalf("pipeline run: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !strings.Contains(r, `"env":"prod"`) {
			t.Errorf("missing env field in: %s", r)
		}
		if !strings.Contains(r, `"app":"logslice"`) {
			t.Errorf("missing app field in: %s", r)
		}
	}
}

func TestAddFieldsFilter_WithPipeline_CombinedWithFieldFilter(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"ok"}`,
		`{"level":"error","msg":"fail"}`,
		`{"level":"info","msg":"also ok"}`,
	}

	ff, err := ParseFieldQuery("level=info")
	if err != nil {
		t.Fatalf("ParseFieldQuery: %v", err)
	}
	af, err := NewAddFieldsFilter([]string{"tagged=yes"})
	if err != nil {
		t.Fatalf("NewAddFieldsFilter: %v", err)
	}

	p := NewPipeline(lines)
	p.AddFilter(ff)
	p.AddTransformer(af)

	results, err := p.Run()
	if err != nil {
		t.Fatalf("pipeline run: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !strings.Contains(r, `"tagged":"yes"`) {
			t.Errorf("missing tagged field in: %s", r)
		}
		if !strings.Contains(r, `"level":"info"`) {
			t.Errorf("expected only info lines, got: %s", r)
		}
	}
}
