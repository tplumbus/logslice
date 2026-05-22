package filter

import (
	"testing"
)

func TestNewRegexFilter_Valid(t *testing.T) {
	f, err := NewRegexFilter("level", `^(error|warn)$`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "level" {
		t.Errorf("expected field 'level', got %q", f.Field())
	}
	if f.Pattern() != `^(error|warn)$` {
		t.Errorf("unexpected pattern: %q", f.Pattern())
	}
}

func TestNewRegexFilter_EmptyField(t *testing.T) {
	_, err := NewRegexFilter("", `.*`)
	if err == nil {
		t.Fatal("expected error for empty field, got nil")
	}
}

func TestNewRegexFilter_EmptyPattern(t *testing.T) {
	_, err := NewRegexFilter("level", "")
	if err == nil {
		t.Fatal("expected error for empty pattern, got nil")
	}
}

func TestNewRegexFilter_InvalidPattern(t *testing.T) {
	_, err := NewRegexFilter("level", `[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestRegexFilter_MatchesLine_Match(t *testing.T) {
	f, _ := NewRegexFilter("level", `^error$`)
	line := `{"level":"error","msg":"something failed"}`
	if !f.MatchesLine(line) {
		t.Error("expected line to match")
	}
}

func TestRegexFilter_MatchesLine_NoMatch(t *testing.T) {
	f, _ := NewRegexFilter("level", `^error$`)
	line := `{"level":"info","msg":"all good"}`
	if f.MatchesLine(line) {
		t.Error("expected line not to match")
	}
}

func TestRegexFilter_MatchesLine_MissingField(t *testing.T) {
	f, _ := NewRegexFilter("level", `.*`)
	line := `{"msg":"no level field here"}`
	if f.MatchesLine(line) {
		t.Error("expected no match when field is absent")
	}
}

func TestRegexFilter_MatchesLine_InvalidJSON(t *testing.T) {
	f, _ := NewRegexFilter("level", `.*`)
	if f.MatchesLine("not-json") {
		t.Error("expected no match for invalid JSON")
	}
}

func TestRegexFilter_MatchesLine_NumericField(t *testing.T) {
	f, _ := NewRegexFilter("status", `^5\d{2}$`)
	line := `{"status":503,"msg":"service unavailable"}`
	if !f.MatchesLine(line) {
		t.Error("expected numeric field to match via string conversion")
	}
}

func TestRegexFilter_MatchesLine_PartialMatch(t *testing.T) {
	f, _ := NewRegexFilter("msg", `timeout`)
	line := `{"level":"error","msg":"connection timeout exceeded"}`
	if !f.MatchesLine(line) {
		t.Error("expected partial regex match to succeed")
	}
}
