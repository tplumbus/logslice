package filter

import (
	"encoding/json"
	"testing"
)

func TestNewFlattenFilter_Valid(t *testing.T) {
	f, err := NewFlattenFilter(".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewFlattenFilter_EmptySep(t *testing.T) {
	_, err := NewFlattenFilter("")
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestFlattenFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewFlattenFilter(".")
	if !f.MatchesLine(`{"a":1}`) {
		t.Error("expected MatchesLine to return true")
	}
	if !f.MatchesLine("not json at all") {
		t.Error("expected MatchesLine to return true for invalid JSON")
	}
}

func TestFlattenFilter_TransformLine_Flat(t *testing.T) {
	f, _ := NewFlattenFilter(".")
	out := f.TransformLine(`{"level":"info","msg":"hello"}`)
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["level"] != "info" || got["msg"] != "hello" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFlattenFilter_TransformLine_Nested(t *testing.T) {
	f, _ := NewFlattenFilter(".")
	out := f.TransformLine(`{"a":{"b":{"c":42}}}`)
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	val, ok := got["a.b.c"]
	if !ok {
		t.Fatalf("expected key \"a.b.c\", got keys: %v", got)
	}
	if val.(float64) != 42 {
		t.Errorf("expected 42, got %v", val)
	}
	if _, nested := got["a"]; nested {
		t.Error("expected nested key \"a\" to be removed")
	}
}

func TestFlattenFilter_TransformLine_CustomSep(t *testing.T) {
	f, _ := NewFlattenFilter("_")
	out := f.TransformLine(`{"x":{"y":"z"}}`)
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := got["x_y"]; !ok {
		t.Errorf("expected key \"x_y\", got: %v", got)
	}
}

func TestFlattenFilter_TransformLine_InvalidJSON(t *testing.T) {
	f, _ := NewFlattenFilter(".")
	input := "not json"
	out := f.TransformLine(input)
	if out != input {
		t.Errorf("expected unchanged line, got: %s", out)
	}
}
