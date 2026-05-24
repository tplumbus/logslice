package filter

import (
	"encoding/json"
	"testing"
)

func TestNewCastFilter_Valid(t *testing.T) {
	for _, tt := range []struct{ field, typ string }{
		{"status", "int"},
		{"score", "float"},
		{"active", "bool"},
		{"code", "string"},
	} {
		_, err := NewCastFilter(tt.field, tt.typ)
		if err != nil {
			t.Errorf("unexpected error for (%s,%s): %v", tt.field, tt.typ, err)
		}
	}
}

func TestNewCastFilter_EmptyField(t *testing.T) {
	_, err := NewCastFilter("", "int")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewCastFilter_InvalidType(t *testing.T) {
	_, err := NewCastFilter("field", "timestamp")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestCastFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewCastFilter("status", "int")
	if !f.MatchesLine(`{"status":"200"}`) {
		t.Fatal("expected MatchesLine to always return true")
	}
}

func TestCastFilter_TransformLine_CastsToInt(t *testing.T) {
	f, _ := NewCastFilter("status", "int")
	out := f.TransformLine(`{"status":"200","msg":"ok"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["status"] != float64(200) {
		t.Errorf("expected 200 (float64), got %v (%T)", obj["status"], obj["status"])
	}
}

func TestCastFilter_TransformLine_CastsToBool(t *testing.T) {
	f, _ := NewCastFilter("active", "bool")
	out := f.TransformLine(`{"active":"true"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["active"] != true {
		t.Errorf("expected true, got %v", obj["active"])
	}
}

func TestCastFilter_TransformLine_MissingField(t *testing.T) {
	f, _ := NewCastFilter("missing", "int")
	original := `{"status":"200"}`
	out := f.TransformLine(original)
	if out != original {
		t.Errorf("expected line unchanged, got %s", out)
	}
}

func TestCastFilter_TransformLine_InvalidJSON(t *testing.T) {
	f, _ := NewCastFilter("status", "int")
	out := f.TransformLine("not-json")
	if out != "not-json" {
		t.Errorf("expected original line on bad JSON, got %s", out)
	}
}

func TestCastFilter_TransformLine_CastToString(t *testing.T) {
	f, _ := NewCastFilter("code", "string")
	out := f.TransformLine(`{"code":404}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["code"] != "404" {
		t.Errorf("expected \"404\", got %v", obj["code"])
	}
}
