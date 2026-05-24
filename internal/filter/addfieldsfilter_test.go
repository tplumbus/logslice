package filter

import (
	"encoding/json"
	"testing"
)

func TestNewAddFieldsFilter_Valid(t *testing.T) {
	f, err := NewAddFieldsFilter([]string{"env=prod", "region=us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewAddFieldsFilter_Empty(t *testing.T) {
	_, err := NewAddFieldsFilter([]string{})
	if err == nil {
		t.Fatal("expected error for empty pairs")
	}
}

func TestNewAddFieldsFilter_InvalidPair(t *testing.T) {
	for _, bad := range []string{"noequals", "=missingkey", ""} {
		_, err := NewAddFieldsFilter([]string{bad})
		if err == nil {
			t.Errorf("expected error for pair %q", bad)
		}
	}
}

func TestAddFieldsFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewAddFieldsFilter([]string{"k=v"})
	if !f.MatchesLine(`{"msg":"hello"}`) {
		t.Error("expected MatchesLine to return true")
	}
	if !f.MatchesLine(`not json at all`) {
		t.Error("expected MatchesLine to return true for non-JSON")
	}
}

func TestAddFieldsFilter_TransformLine_AddsField(t *testing.T) {
	f, _ := NewAddFieldsFilter([]string{"env=staging"})
	out := f.TransformLine(`{"msg":"hello"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if obj["env"] != "staging" {
		t.Errorf("expected env=staging, got %v", obj["env"])
	}
	if obj["msg"] != "hello" {
		t.Errorf("expected original msg preserved, got %v", obj["msg"])
	}
}

func TestAddFieldsFilter_TransformLine_OverwritesField(t *testing.T) {
	f, _ := NewAddFieldsFilter([]string{"level=info"})
	out := f.TransformLine(`{"level":"debug","msg":"test"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if obj["level"] != "info" {
		t.Errorf("expected level overwritten to info, got %v", obj["level"])
	}
}

func TestAddFieldsFilter_TransformLine_InvalidJSON(t *testing.T) {
	f, _ := NewAddFieldsFilter([]string{"env=prod"})
	raw := `not json`
	out := f.TransformLine(raw)
	if out != raw {
		t.Errorf("expected original line returned for invalid JSON, got %q", out)
	}
}
