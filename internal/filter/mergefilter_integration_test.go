package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestMergeFilter_WithPipeline_AddsFields(t *testing.T) {
	lines := []string{
		`{"msg":"startup","level":"info"}`,
		`{"msg":"shutdown","level":"warn"}`,
	}
	input := strings.NewReader(strings.Join(lines, "\n") + "\n")
	var buf bytes.Buffer

	f, err := NewMergeFilter([]string{"env=staging", "app=logslice"}, false)
	if err != nil {
		t.Fatalf("NewMergeFilter: %v", err)
	}

	p := NewPipeline(input, &buf)
	p.AddFilter(f)
	if err := p.Run(); err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	outLines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(outLines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(outLines))
	}
	for _, l := range outLines {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(l), &obj); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if obj["env"] != "staging" {
			t.Errorf("expected env=staging, got %v", obj["env"])
		}
		if obj["app"] != "logslice" {
			t.Errorf("expected app=logslice, got %v", obj["app"])
		}
	}
}

func TestMergeFilter_WithPipeline_OverwriteMode(t *testing.T) {
	lines := []string{
		`{"msg":"hello","env":"dev"}`,
	}
	input := strings.NewReader(strings.Join(lines, "\n") + "\n")
	var buf bytes.Buffer

	f, _ := NewMergeFilter([]string{"env=prod"}, true)
	p := NewPipeline(input, &buf)
	p.AddFilter(f)
	if err := p.Run(); err != nil {
		t.Fatalf("pipeline run: %v", err)
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("expected env=prod after overwrite, got %v", obj["env"])
	}
}
