package filter

import (
	"encoding/json"
	"testing"
)

func TestNewTrimFilter_Valid(t *testing.T) {
	f, err := NewTrimFilter([]string{"msg"}, "both")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewTrimFilter_EmptyFields(t *testing.T) {
	_, err := NewTrimFilter([]string{}, "both")
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNewTrimFilter_InvalidMode(t *testing.T) {
	_, err := NewTrimFilter([]string{"msg"}, "center")
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestNewTrimFilter_BlankFieldsFiltered(t *testing.T) {
	_, err := NewTrimFilter([]string{"  ", ""}, "both")
	if err == nil {
		t.Fatal("expected error when all fields are blank")
	}
}

func TestTrimFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewTrimFilter([]string{"msg"}, "both")
	if !f.MatchesLine(`{"msg":"hello"}`) {
		t.Fatal("expected MatchesLine to return true")
	}
	if !f.MatchesLine(`not json`) {
		t.Fatal("expected MatchesLine to return true for invalid JSON")
	}
}

func TestTrimFilter_TransformLine_TrimsBoth(t *testing.T) {
	f, _ := NewTrimFilter([]string{"msg"}, "both")
	out, err := f.TransformLine(`{"msg":"  hello world  "}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if obj["msg"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", obj["msg"])
	}
}

func TestTrimFilter_TransformLine_TrimsLeft(t *testing.T) {
	f, _ := NewTrimFilter([]string{"msg"}, "left")
	out, _ := f.TransformLine(`{"msg":"  hello  "}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["msg"] != "hello  " {
		t.Errorf("expected 'hello  ', got %q", obj["msg"])
	}
}

func TestTrimFilter_TransformLine_TrimsRight(t *testing.T) {
	f, _ := NewTrimFilter([]string{"msg"}, "right")
	out, _ := f.TransformLine(`{"msg":"  hello  "}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["msg"] != "  hello" {
		t.Errorf("expected '  hello', got %q", obj["msg"])
	}
}

func TestTrimFilter_TransformLine_NonStringFieldSkipped(t *testing.T) {
	f, _ := NewTrimFilter([]string{"count"}, "both")
	original := `{"count":42}`
	out, _ := f.TransformLine(original)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if v, ok := obj["count"].(float64); !ok || v != 42 {
		t.Errorf("expected count=42, got %v", obj["count"])
	}
}

func TestParseTrimFlag_Valid(t *testing.T) {
	f, err := ParseTrimFlag("msg,user:both")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(f.fields))
	}
}

func TestParseTrimFlag_MissingColon(t *testing.T) {
	_, err := ParseTrimFlag("msgboth")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}
