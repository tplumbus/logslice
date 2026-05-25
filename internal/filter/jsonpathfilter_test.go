package filter

import (
	"testing"
)

func TestNewJSONPathFilter_Valid(t *testing.T) {
	f, err := NewJSONPathFilter("meta.user.id", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewJSONPathFilter_EmptyPath(t *testing.T) {
	_, err := NewJSONPathFilter("", "val")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewJSONPathFilter_EmptySegment(t *testing.T) {
	_, err := NewJSONPathFilter("meta..id", "val")
	if err == nil {
		t.Fatal("expected error for empty segment")
	}
}

func TestJSONPathFilter_MatchesLine_TopLevel(t *testing.T) {
	f, _ := NewJSONPathFilter("status", "ok")
	line := `{"status":"ok","code":200}`
	if !f.MatchesLine(line) {
		t.Error("expected match")
	}
}

func TestJSONPathFilter_MatchesLine_Nested(t *testing.T) {
	f, _ := NewJSONPathFilter("meta.user.id", "42")
	line := `{"meta":{"user":{"id":"42","name":"alice"}}}`
	if !f.MatchesLine(line) {
		t.Error("expected match on nested path")
	}
}

func TestJSONPathFilter_MatchesLine_NoMatch(t *testing.T) {
	f, _ := NewJSONPathFilter("meta.user.id", "99")
	line := `{"meta":{"user":{"id":"42"}}}`
	if f.MatchesLine(line) {
		t.Error("expected no match")
	}
}

func TestJSONPathFilter_MatchesLine_MissingPath(t *testing.T) {
	f, _ := NewJSONPathFilter("meta.region", "us-east")
	line := `{"meta":{"user":{"id":"1"}}}`
	if f.MatchesLine(line) {
		t.Error("expected no match for missing path")
	}
}

func TestJSONPathFilter_MatchesLine_InvalidJSON(t *testing.T) {
	f, _ := NewJSONPathFilter("a.b", "v")
	if f.MatchesLine("not-json") {
		t.Error("expected no match for invalid JSON")
	}
}

func TestJSONPathFilter_TransformLine_Unchanged(t *testing.T) {
	f, _ := NewJSONPathFilter("a", "b")
	line := `{"a":"b"}`
	if got := f.TransformLine(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestParseJSONPathFlag_Valid(t *testing.T) {
	f, err := ParseJSONPathFlag("meta.env=production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"meta":{"env":"production"}}`
	if !f.MatchesLine(line) {
		t.Error("expected match after parse")
	}
}

func TestParseJSONPathFlag_MissingEquals(t *testing.T) {
	_, err := ParseJSONPathFlag("meta.env")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseJSONPathFlag_EmptyKey(t *testing.T) {
	_, err := ParseJSONPathFlag("=value")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}
