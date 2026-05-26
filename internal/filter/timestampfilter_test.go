package filter

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewTimestampFilter_Valid(t *testing.T) {
	f, err := NewTimestampFilter("ts", "2006-01-02", "01/02/2006")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewTimestampFilter_EmptyField(t *testing.T) {
	_, err := NewTimestampFilter("", "2006-01-02", "01/02/2006")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTimestampFilter_EmptyInFmt(t *testing.T) {
	_, err := NewTimestampFilter("ts", "", "01/02/2006")
	if err == nil {
		t.Fatal("expected error for empty inFmt")
	}
}

func TestNewTimestampFilter_EmptyOutFmt(t *testing.T) {
	_, err := NewTimestampFilter("ts", "2006-01-02", "")
	if err == nil {
		t.Fatal("expected error for empty outFmt")
	}
}

func TestTimestampFilter_MatchesLine_AlwaysTrue(t *testing.T) {
	f, _ := NewTimestampFilter("ts", "2006-01-02", "01/02/2006")
	if !f.MatchesLine(`{"ts":"2024-01-15"}`) {
		t.Fatal("expected MatchesLine to return true")
	}
}

func TestTimestampFilter_TransformLine_Reformats(t *testing.T) {
	f, _ := NewTimestampFilter("ts", "2006-01-02", "01/02/2006")
	out := f.TransformLine(`{"ts":"2024-03-07","level":"info"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if obj["ts"] != "03/07/2024" {
		t.Errorf("expected 03/07/2024, got %v", obj["ts"])
	}
}

func TestTimestampFilter_TransformLine_MissingField(t *testing.T) {
	f, _ := NewTimestampFilter("ts", "2006-01-02", "01/02/2006")
	line := `{"level":"info"}`
	if out := f.TransformLine(line); out != line {
		t.Errorf("expected unchanged line, got %s", out)
	}
}

func TestTimestampFilter_TransformLine_UnparseableValue(t *testing.T) {
	f, _ := NewTimestampFilter("ts", "2006-01-02", "01/02/2006")
	line := `{"ts":"not-a-date"}`
	if out := f.TransformLine(line); out != line {
		t.Errorf("expected unchanged line, got %s", out)
	}
}

func TestParseTimestampFlag_Valid(t *testing.T) {
	f, err := ParseTimestampFlag("ts:rfc3339:date")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.inFmt != timestampAliases["rfc3339"] {
		t.Errorf("alias not resolved for inFmt")
	}
	if f.outFmt != timestampAliases["date"] {
		t.Errorf("alias not resolved for outFmt")
	}
}

func TestParseTimestampFlag_MissingParts(t *testing.T) {
	_, err := ParseTimestampFlag("ts:rfc3339")
	if err == nil {
		t.Fatal("expected error for missing out-format")
	}
	if !strings.Contains(err.Error(), "expected field:in-format:out-format") {
		t.Errorf("unexpected error message: %v", err)
	}
}
