package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestNewDedupFilter_Valid(t *testing.T) {
	_, err := filter.NewDedupFilter("request_id", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewDedupFilter_EmptyField(t *testing.T) {
	_, err := filter.NewDedupFilter("", 0)
	if err == nil {
		t.Fatal("expected error for empty field, got nil")
	}
}

func TestDedupFilter_MatchesLine_UniqueLines(t *testing.T) {
	f, _ := filter.NewDedupFilter("id", 0)

	lines := []string{
		`{"id":"aaa","msg":"first"}`,
		`{"id":"bbb","msg":"second"}`,
		`{"id":"ccc","msg":"third"}`,
	}

	for _, line := range lines {
		if !f.MatchesLine(line) {
			t.Errorf("expected line to pass dedup filter: %s", line)
		}
	}
}

func TestDedupFilter_MatchesLine_DuplicateBlocked(t *testing.T) {
	f, _ := filter.NewDedupFilter("id", 0)

	line := `{"id":"aaa","msg":"first"}`
	if !f.MatchesLine(line) {
		t.Fatal("first occurrence should pass")
	}
	if f.MatchesLine(line) {
		t.Fatal("duplicate occurrence should be blocked")
	}
}

func TestDedupFilter_MatchesLine_MissingField(t *testing.T) {
	f, _ := filter.NewDedupFilter("id", 0)

	line := `{"msg":"no id here"}`
	if !f.MatchesLine(line) {
		t.Fatal("line missing dedup field should pass through")
	}
}

func TestDedupFilter_MatchesLine_InvalidJSON(t *testing.T) {
	f, _ := filter.NewDedupFilter("id", 0)

	if !f.MatchesLine("not json at all") {
		t.Fatal("invalid JSON should pass through")
	}
}

func TestDedupFilter_MatchesLine_MaxSeen(t *testing.T) {
	// maxSeen=2 means each unique value passes at most 2 times
	f, _ := filter.NewDedupFilter("level", 2)

	line := `{"level":"error","msg":"boom"}`
	if !f.MatchesLine(line) {
		t.Fatal("first occurrence should pass")
	}
	if !f.MatchesLine(line) {
		t.Fatal("second occurrence should pass (maxSeen=2)")
	}
	if f.MatchesLine(line) {
		t.Fatal("third occurrence should be blocked")
	}
}

func TestDedupFilter_MatchesLine_DifferentValues(t *testing.T) {
	f, _ := filter.NewDedupFilter("host", 0)

	if !f.MatchesLine(`{"host":"web-1"}`) {
		t.Fatal("first host should pass")
	}
	if !f.MatchesLine(`{"host":"web-2"}`) {
		t.Fatal("different host should pass")
	}
	if f.MatchesLine(`{"host":"web-1"}`) {
		t.Fatal("repeated host should be blocked")
	}
}
