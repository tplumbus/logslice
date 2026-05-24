package filter

import (
	"encoding/json"
	"testing"
)

func TestNewSelectFilter_Valid(t *testing.T) {
	sf, err := NewSelectFilter([]string{"level", "msg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sf == nil {
		t.Fatal("expected non-nil SelectFilter")
	}
}

func TestNewSelectFilter_EmptyFields(t *testing.T) {
	_, err := NewSelectFilter([]string{"", "  "})
	if err == nil {
		t.Fatal("expected error for empty field list")
	}
}

func TestNewSelectFilter_SkipsBlanks(t *testing.T) {
	sf, err := NewSelectFilter([]string{"", "level", "  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sf.fields) != 1 || sf.fields[0] != "level" {
		t.Fatalf("expected only 'level', got %v", sf.fields)
	}
}

func TestSelectFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	sf, _ := NewSelectFilter([]string{"msg"})
	if !sf.MatchesLine(`{"msg":"hello","level":"info"}`) {
		t.Fatal("expected MatchesLine to return true")
	}
	if !sf.MatchesLine(`not json at all`) {
		t.Fatal("expected MatchesLine to return true for invalid JSON")
	}
}

func TestSelectFilter_TransformLine_KeepsFields(t *testing.T) {
	sf, _ := NewSelectFilter([]string{"level", "msg"})
	input := `{"ts":"2024-01-01","level":"error","msg":"oops","code":500}`
	out := sf.TransformLine(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := m["ts"]; ok {
		t.Error("field 'ts' should have been removed")
	}
	if _, ok := m["code"]; ok {
		t.Error("field 'code' should have been removed")
	}
	if m["level"] != "error" {
		t.Errorf("expected level=error, got %v", m["level"])
	}
	if m["msg"] != "oops" {
		t.Errorf("expected msg=oops, got %v", m["msg"])
	}
}

func TestSelectFilter_TransformLine_MissingFieldOmitted(t *testing.T) {
	sf, _ := NewSelectFilter([]string{"level", "trace_id"})
	input := `{"level":"info","msg":"hi"}`
	out := sf.TransformLine(input)

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := m["trace_id"]; ok {
		t.Error("missing field 'trace_id' should not appear in output")
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}

func TestSelectFilter_TransformLine_InvalidJSON(t *testing.T) {
	sf, _ := NewSelectFilter([]string{"msg"})
	input := `not json`
	out := sf.TransformLine(input)
	if out != input {
		t.Errorf("expected original line returned for invalid JSON, got %q", out)
	}
}

func TestParseSelectFields_Valid(t *testing.T) {
	sf, err := ParseSelectFields("level, msg, ts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sf.fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(sf.fields))
	}
}

func TestParseSelectFields_Empty(t *testing.T) {
	_, err := ParseSelectFields("   ")
	if err == nil {
		t.Fatal("expected error for blank input")
	}
}
