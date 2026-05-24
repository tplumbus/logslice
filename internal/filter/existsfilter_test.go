package filter

import (
	"strings"
	"testing"
)

func TestNewExistsFilter_Valid(t *testing.T) {
	f, err := NewExistsFilter("level", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewExistsFilter_EmptyField(t *testing.T) {
	_, err := NewExistsFilter("", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestExistsFilter_MatchesLine_FieldPresent(t *testing.T) {
	f, _ := NewExistsFilter("level", false)
	line := `{"level":"info","msg":"hello"}`
	if !f.MatchesLine(line) {
		t.Errorf("expected match when field present")
	}
}

func TestExistsFilter_MatchesLine_FieldAbsent(t *testing.T) {
	f, _ := NewExistsFilter("level", false)
	line := `{"msg":"hello"}`
	if f.MatchesLine(line) {
		t.Errorf("expected no match when field absent")
	}
}

func TestExistsFilter_MatchesLine_NegateAbsent(t *testing.T) {
	f, _ := NewExistsFilter("level", true)
	line := `{"msg":"hello"}`
	if !f.MatchesLine(line) {
		t.Errorf("expected match when field absent and negate=true")
	}
}

func TestExistsFilter_MatchesLine_NegatePresent(t *testing.T) {
	f, _ := NewExistsFilter("level", true)
	line := `{"level":"info","msg":"hello"}`
	if f.MatchesLine(line) {
		t.Errorf("expected no match when field present and negate=true")
	}
}

func TestExistsFilter_MatchesLine_InvalidJSON(t *testing.T) {
	f, _ := NewExistsFilter("level", false)
	if f.MatchesLine("not json") {
		t.Errorf("expected false for invalid JSON")
	}
}

func TestExistsFilter_TransformLine_Unchanged(t *testing.T) {
	f, _ := NewExistsFilter("level", false)
	line := `{"level":"info"}`
	if got := f.TransformLine(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestExistsFilter_WithPipeline_FiltersPresent(t *testing.T) {
	input := strings.NewReader(
		`{"level":"info","msg":"a"}` + "\n" +
			`{"msg":"b"}` + "\n" +
			`{"level":"warn","msg":"c"}` + "\n",
	)
	f, _ := NewExistsFilter("level", false)
	p := NewPipeline(input, f)
	var out []string
	_ = p.Run(func(line string) error {
		out = append(out, line)
		return nil
	})
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(out), out)
	}
}

func TestParseExistsFlag_Valid(t *testing.T) {
	f, err := ParseExistsFlag("status", false)
	if err != nil || f == nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseExistsFlag_Empty(t *testing.T) {
	_, err := ParseExistsFlag("", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
