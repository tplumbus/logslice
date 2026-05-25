package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestSplitFilter_WithPipeline_SplitsField(t *testing.T) {
	lines := []string{
		`{"level":"info","tags":"auth,login,success"}`,
		`{"level":"warn","tags":"db,timeout"}`,
		`{"level":"error","tags":"crash"}`,
	}

	f, err := NewSplitFilter("tags", ",", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := NewPipeline(strings.NewReader(strings.Join(lines, "\n")))
	p.AddTransformer(f)

	var buf bytes.Buffer
	if err := p.Run(&buf); err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	results := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(results) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(results))
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(results[0]), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	arr, ok := obj["tags"].([]interface{})
	if !ok || len(arr) != 3 {
		t.Errorf("expected 3 tags, got %v", obj["tags"])
	}
}

func TestSplitFilter_WithPipeline_CombinedWithFieldFilter(t *testing.T) {
	lines := []string{
		`{"level":"info","tags":"auth,login"}`,
		`{"level":"warn","tags":"db,timeout"}`,
	}

	ff, err := ParseFieldQuery("level=info")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sf, err := NewSplitFilter("tags", ",", "tag_list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := NewPipeline(strings.NewReader(strings.Join(lines, "\n")))
	p.AddFilter(ff)
	p.AddTransformer(sf)

	var buf bytes.Buffer
	if err := p.Run(&buf); err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	results := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(results) != 1 {
		t.Fatalf("expected 1 line, got %d", len(results))
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(results[0]), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := obj["tag_list"]; !ok {
		t.Error("expected 'tag_list' field in output")
	}
	if obj["level"] != "info" {
		t.Errorf("expected level=info, got %v", obj["level"])
	}
}
