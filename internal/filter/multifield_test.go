package filter

import (
	"testing"
)

func TestNewMultiFieldFilter_Valid(t *testing.T) {
	mf, err := NewMultiFieldFilter([]string{"level=error", "service=api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mf.Len() != 2 {
		t.Errorf("expected 2 queries, got %d", mf.Len())
	}
}

func TestNewMultiFieldFilter_SkipsBlanks(t *testing.T) {
	mf, err := NewMultiFieldFilter([]string{"level=info", "", "  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mf.Len() != 1 {
		t.Errorf("expected 1 query, got %d", mf.Len())
	}
}

func TestNewMultiFieldFilter_InvalidQuery(t *testing.T) {
	_, err := NewMultiFieldFilter([]string{"badquery"})
	if err == nil {
		t.Fatal("expected error for invalid query, got nil")
	}
}

func TestMultiFieldFilter_MatchesLine_AllMatch(t *testing.T) {
	mf, _ := NewMultiFieldFilter([]string{"level=error", "service=api"})
	line := `{"level":"error","service":"api","msg":"oops"}`
	if !mf.MatchesLine(line) {
		t.Error("expected line to match all field queries")
	}
}

func TestMultiFieldFilter_MatchesLine_PartialMatch(t *testing.T) {
	mf, _ := NewMultiFieldFilter([]string{"level=error", "service=api"})
	line := `{"level":"error","service":"worker","msg":"oops"}`
	if mf.MatchesLine(line) {
		t.Error("expected line NOT to match when one field differs")
	}
}

func TestMultiFieldFilter_MatchesLine_EmptyFilter(t *testing.T) {
	mf, _ := NewMultiFieldFilter([]string{})
	line := `{"level":"debug","msg":"anything"}`
	if !mf.MatchesLine(line) {
		t.Error("empty filter should match every line")
	}
}

func TestMultiFieldFilter_MatchesLine_InvalidJSON(t *testing.T) {
	mf, _ := NewMultiFieldFilter([]string{"level=error"})
	if mf.MatchesLine("not json at all") {
		t.Error("invalid JSON should not match")
	}
}

func TestMultiFieldFilter_MatchesLine_MissingKey(t *testing.T) {
	mf, _ := NewMultiFieldFilter([]string{"level=error"})
	line := `{"msg":"no level field here"}`
	if mf.MatchesLine(line) {
		t.Error("line missing the queried key should not match")
	}
}
