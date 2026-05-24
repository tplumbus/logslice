package filter

import (
	"strings"
	"testing"
)

func TestFlattenFilter_WithPipeline_FlattensNestedLogs(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:00Z","http":{"method":"GET","status":200}}`,
		`{"ts":"2024-01-01T00:01:00Z","http":{"method":"POST","status":201}}`,
		`not json`,
	}

	ff, err := NewFlattenFilter(".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := NewPipeline(lines, nil, []Filter{ff})
	results, err := p.Run()
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	// Two valid JSON lines + one passthrough
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	for _, r := range results[:2] {
		if !strings.Contains(r, "http.method") {
			t.Errorf("expected flattened key \"http.method\" in: %s", r)
		}
		if !strings.Contains(r, "http.status") {
			t.Errorf("expected flattened key \"http.status\" in: %s", r)
		}
		if strings.Contains(r, `"http":{`) {
			t.Errorf("expected nested \"http\" to be removed in: %s", r)
		}
	}

	if results[2] != "not json" {
		t.Errorf("expected passthrough for invalid JSON, got: %s", results[2])
	}
}

func TestFlattenFilter_WithPipeline_CombinedWithFieldFilter(t *testing.T) {
	lines := []string{
		`{"user":{"id":"u1","role":"admin"}}`,
		`{"user":{"id":"u2","role":"viewer"}}`,
	}

	ff, _ := NewFlattenFilter(".")
	fq, err := ParseFieldQuery("user.role=admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := NewPipeline(lines, nil, []Filter{ff, fq})
	results, err := p.Run()
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d: %v", len(results), results)
	}
	if !strings.Contains(results[0], "admin") {
		t.Errorf("expected admin line, got: %s", results[0])
	}
}
