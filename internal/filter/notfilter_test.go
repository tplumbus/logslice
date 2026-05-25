package filter

import (
	"testing"
)

func TestNewNotFilter_Valid(t *testing.T) {
	fq, _ := ParseFieldQuery("level=debug")
	_, err := NewNotFilter(fq)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewNotFilter_NoFilters(t *testing.T) {
	_, err := NewNotFilter()
	if err == nil {
		t.Fatal("expected error for empty filters")
	}
}

func TestNewNotFilter_NilFilter(t *testing.T) {
	_, err := NewNotFilter(nil)
	if err == nil {
		t.Fatal("expected error for nil filter")
	}
}

func TestNotFilter_MatchesLine_NoneMatch(t *testing.T) {
	fq, _ := ParseFieldQuery("level=debug")
	nf, _ := NewNotFilter(fq)
	line := `{"level":"info","msg":"ok"}`
	if !nf.MatchesLine(line) {
		t.Error("expected line to pass (none match)")
	}
}

func TestNotFilter_MatchesLine_OneMatches(t *testing.T) {
	fq, _ := ParseFieldQuery("level=debug")
	nf, _ := NewNotFilter(fq)
	line := `{"level":"debug","msg":"verbose"}`
	if nf.MatchesLine(line) {
		t.Error("expected line to be rejected (one matches)")
	}
}

func TestNotFilter_MatchesLine_MultipleFilters(t *testing.T) {
	fq1, _ := ParseFieldQuery("level=debug")
	fq2, _ := ParseFieldQuery("env=staging")
	nf, _ := NewNotFilter(fq1, fq2)

	matches := []struct {
		line string
		want bool
	}{
		{`{"level":"info","env":"prod"}`, true},
		{`{"level":"debug","env":"prod"}`, false},
		{`{"level":"info","env":"staging"}`, false},
		{`{"level":"debug","env":"staging"}`, false},
	}
	for _, tc := range matches {
		got := nf.MatchesLine(tc.line)
		if got != tc.want {
			t.Errorf("line %s: got %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestParseNotFlag_Valid(t *testing.T) {
	nf, err := ParseNotFlag("level=debug,env=staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"level":"info","env":"prod"}`
	if !nf.MatchesLine(line) {
		t.Error("expected line to pass")
	}
}

func TestParseNotFlag_Empty(t *testing.T) {
	_, err := ParseNotFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}

func TestParseNotFlag_InvalidQuery(t *testing.T) {
	_, err := ParseNotFlag("nodEquals")
	if err == nil {
		t.Fatal("expected error for invalid query")
	}
}
