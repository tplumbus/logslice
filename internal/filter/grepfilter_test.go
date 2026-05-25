package filter

import (
	"encoding/json"
	"testing"
)

func makeGrepLine(t *testing.T, fields map[string]interface{}) string {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}

func TestNewGrepFilter_Valid(t *testing.T) {
	f, err := NewGrepFilter("msg", "error", false)
	if err != nil || f == nil {
		t.Fatalf("expected valid filter, got err=%v", err)
	}
}

func TestNewGrepFilter_EmptyField(t *testing.T) {
	_, err := NewGrepFilter("", "error", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewGrepFilter_EmptySubstring(t *testing.T) {
	_, err := NewGrepFilter("msg", "", false)
	if err == nil {
		t.Fatal("expected error for empty substring")
	}
}

func TestGrepFilter_MatchesLine_Contains(t *testing.T) {
	f, _ := NewGrepFilter("msg", "fail", false)
	line := makeGrepLine(t, map[string]interface{}{"msg": "request fail"})
	if !f.MatchesLine(line) {
		t.Error("expected match")
	}
}

func TestGrepFilter_MatchesLine_NoMatch(t *testing.T) {
	f, _ := NewGrepFilter("msg", "error", false)
	line := makeGrepLine(t, map[string]interface{}{"msg": "all good"})
	if f.MatchesLine(line) {
		t.Error("expected no match")
	}
}

func TestGrepFilter_MatchesLine_CaseSensitive(t *testing.T) {
	f, _ := NewGrepFilter("msg", "Error", false)
	line := makeGrepLine(t, map[string]interface{}{"msg": "error occurred"})
	if f.MatchesLine(line) {
		t.Error("expected no match (case-sensitive)")
	}
}

func TestGrepFilter_MatchesLine_IgnoreCase(t *testing.T) {
	f, _ := NewGrepFilter("msg", "error", true)
	line := makeGrepLine(t, map[string]interface{}{"msg": "Error occurred"})
	if !f.MatchesLine(line) {
		t.Error("expected match (case-insensitive)")
	}
}

func TestGrepFilter_MatchesLine_MissingField(t *testing.T) {
	f, _ := NewGrepFilter("msg", "error", false)
	line := makeGrepLine(t, map[string]interface{}{"level": "warn"})
	if f.MatchesLine(line) {
		t.Error("expected no match for missing field")
	}
}

func TestGrepFilter_MatchesLine_InvalidJSON(t *testing.T) {
	f, _ := NewGrepFilter("msg", "error", false)
	if f.MatchesLine("not-json") {
		t.Error("expected no match for invalid JSON")
	}
}

func TestParseGrepFlag_Valid(t *testing.T) {
	f, err := ParseGrepFlag("msg:error")
	if err != nil || f == nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.ignoreCase {
		t.Error("expected case-sensitive")
	}
}

func TestParseGrepFlag_IgnoreCase(t *testing.T) {
	f, err := ParseGrepFlag("msg:error:i")
	if err != nil || f == nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.ignoreCase {
		t.Error("expected case-insensitive")
	}
}

func TestParseGrepFlag_MissingParts(t *testing.T) {
	_, err := ParseGrepFlag("msgonly")
	if err == nil {
		t.Fatal("expected error for missing parts")
	}
}
