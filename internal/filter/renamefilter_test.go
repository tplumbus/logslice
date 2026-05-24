package filter

import (
	"encoding/json"
	"testing"
)

func TestNewRenameFilter_Valid(t *testing.T) {
	f, err := NewRenameFilter([]string{"msg=message", "ts=timestamp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewRenameFilter_Empty(t *testing.T) {
	_, err := NewRenameFilter([]string{})
	if err == nil {
		t.Fatal("expected error for empty pairs")
	}
}

func TestNewRenameFilter_InvalidPair(t *testing.T) {
	cases := []string{"noequals", "=newonly", "oldonly="}
	for _, c := range cases {
		_, err := NewRenameFilter([]string{c})
		if err == nil {
			t.Errorf("expected error for pair %q", c)
		}
	}
}

func TestRenameFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewRenameFilter([]string{"a=b"})
	if !f.MatchesLine(`{"a":"v"}`) {
		t.Error("expected MatchesLine to always return true")
	}
	if !f.MatchesLine(`not json at all`) {
		t.Error("expected MatchesLine to return true for invalid JSON")
	}
}

func TestRenameFilter_TransformLine_RenamesField(t *testing.T) {
	f, _ := NewRenameFilter([]string{"msg=message"})
	out := f.TransformLine(`{"msg":"hello","level":"info"}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := obj["msg"]; ok {
		t.Error("old field 'msg' should have been removed")
	}
	if obj["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", obj["message"])
	}
	if obj["level"] != "info" {
		t.Errorf("expected level=info to be preserved, got %v", obj["level"])
	}
}

func TestRenameFilter_TransformLine_MissingFieldSkipped(t *testing.T) {
	f, _ := NewRenameFilter([]string{"missing=new"})
	out := f.TransformLine(`{"level":"warn"}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := obj["new"]; ok {
		t.Error("field 'new' should not exist when source was missing")
	}
}

func TestRenameFilter_TransformLine_InvalidJSON(t *testing.T) {
	f, _ := NewRenameFilter([]string{"a=b"})
	input := `not valid json`
	out := f.TransformLine(input)
	if out != input {
		t.Errorf("expected original line returned for invalid JSON, got %q", out)
	}
}

func TestRenameFilter_TransformLine_MultipleRenames(t *testing.T) {
	f, _ := NewRenameFilter([]string{"ts=timestamp", "lvl=level"})
	out := f.TransformLine(`{"ts":"2024-01-01","lvl":"error","msg":"oops"}`)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := obj["ts"]; ok {
		t.Error("old field 'ts' should be removed")
	}
	if _, ok := obj["lvl"]; ok {
		t.Error("old field 'lvl' should be removed")
	}
	if obj["timestamp"] != "2024-01-01" {
		t.Errorf("expected timestamp=2024-01-01, got %v", obj["timestamp"])
	}
	if obj["level"] != "error" {
		t.Errorf("expected level=error, got %v", obj["level"])
	}
}
