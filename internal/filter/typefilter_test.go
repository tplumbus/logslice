package filter

import (
	"testing"
)

func TestNewTypeFilter_Valid(t *testing.T) {
	for _, typ := range []string{"string", "number", "bool", "null", "array", "object"} {
		_, err := NewTypeFilter("field", typ)
		if err != nil {
			t.Errorf("expected no error for type %q, got %v", typ, err)
		}
	}
}

func TestNewTypeFilter_EmptyField(t *testing.T) {
	_, err := NewTypeFilter("", "string")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTypeFilter_InvalidType(t *testing.T) {
	_, err := NewTypeFilter("level", "integer")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestTypeFilter_MatchesLine_String(t *testing.T) {
	f, _ := NewTypeFilter("level", "string")
	if !f.MatchesLine(`{"level":"info","msg":"ok"}`) {
		t.Error("expected match for string field")
	}
	if f.MatchesLine(`{"level":42}`) {
		t.Error("expected no match: level is number not string")
	}
}

func TestTypeFilter_MatchesLine_Number(t *testing.T) {
	f, _ := NewTypeFilter("status", "number")
	if !f.MatchesLine(`{"status":200}`) {
		t.Error("expected match for numeric field")
	}
	if f.MatchesLine(`{"status":"200"}`) {
		t.Error("expected no match: status is string")
	}
}

func TestTypeFilter_MatchesLine_Bool(t *testing.T) {
	f, _ := NewTypeFilter("ok", "bool")
	if !f.MatchesLine(`{"ok":true}`) {
		t.Error("expected match for bool true")
	}
	if !f.MatchesLine(`{"ok":false}`) {
		t.Error("expected match for bool false")
	}
}

func TestTypeFilter_MatchesLine_Null(t *testing.T) {
	f, _ := NewTypeFilter("err", "null")
	if !f.MatchesLine(`{"err":null}`) {
		t.Error("expected match for null")
	}
}

func TestTypeFilter_MatchesLine_Array(t *testing.T) {
	f, _ := NewTypeFilter("tags", "array")
	if !f.MatchesLine(`{"tags":["a","b"]}`) {
		t.Error("expected match for array")
	}
}

func TestTypeFilter_MatchesLine_Object(t *testing.T) {
	f, _ := NewTypeFilter("meta", "object")
	if !f.MatchesLine(`{"meta":{"k":"v"}}`) {
		t.Error("expected match for object")
	}
}

func TestTypeFilter_MatchesLine_MissingField(t *testing.T) {
	f, _ := NewTypeFilter("missing", "string")
	if f.MatchesLine(`{"level":"info"}`) {
		t.Error("expected no match when field is absent")
	}
}

func TestTypeFilter_MatchesLine_InvalidJSON(t *testing.T) {
	f, _ := NewTypeFilter("level", "string")
	if f.MatchesLine(`not json`) {
		t.Error("expected no match for invalid JSON")
	}
}

func TestParseTypeFlag_Valid(t *testing.T) {
	f, err := ParseTypeFlag("level:string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.field != "level" || f.wantType != "string" {
		t.Errorf("unexpected parsed values: %+v", f)
	}
}

func TestParseTypeFlag_MissingColon(t *testing.T) {
	_, err := ParseTypeFlag("levelstring")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}
