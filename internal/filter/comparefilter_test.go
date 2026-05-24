package filter

import (
	"fmt"
	"testing"
)

func makeCompareLine(field string, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf(`{%q:%q}`, field, v)
	case float64:
		return fmt.Sprintf(`{%q:%v}`, field, v)
	case int:
		return fmt.Sprintf(`{%q:%d}`, field, v)
	}
	return `{}`
}

func TestNewCompareFilter_Valid(t *testing.T) {
	f, err := NewCompareFilter("latency", OpGte, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewCompareFilter_EmptyField(t *testing.T) {
	_, err := NewCompareFilter("", OpEq, 1)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewCompareFilter_InvalidOp(t *testing.T) {
	_, err := NewCompareFilter("x", "between", 5)
	if err == nil {
		t.Fatal("expected error for unknown operator")
	}
}

func TestCompareFilter_MatchesLine_Gt(t *testing.T) {
	f, _ := NewCompareFilter("score", OpGt, 50)
	if !f.MatchesLine(makeCompareLine("score", float64(99))) {
		t.Error("expected match for score=99 > 50")
	}
	if f.MatchesLine(makeCompareLine("score", float64(50))) {
		t.Error("expected no match for score=50 (not strictly greater)")
	}
}

func TestCompareFilter_MatchesLine_Eq(t *testing.T) {
	f, _ := NewCompareFilter("code", OpEq, 200)
	if !f.MatchesLine(makeCompareLine("code", float64(200))) {
		t.Error("expected match for code=200")
	}
	if f.MatchesLine(makeCompareLine("code", float64(404))) {
		t.Error("expected no match for code=404")
	}
}

func TestCompareFilter_MatchesLine_StringNumeric(t *testing.T) {
	f, _ := NewCompareFilter("latency", OpLt, 100)
	if !f.MatchesLine(makeCompareLine("latency", "42")) {
		t.Error("expected match for string numeric latency=42 < 100")
	}
}

func TestCompareFilter_MatchesLine_MissingField(t *testing.T) {
	f, _ := NewCompareFilter("missing", OpGt, 0)
	if f.MatchesLine(`{"other":1}`) {
		t.Error("expected no match for missing field")
	}
}

func TestCompareFilter_MatchesLine_InvalidJSON(t *testing.T) {
	f, _ := NewCompareFilter("x", OpGt, 0)
	if f.MatchesLine(`not json`) {
		t.Error("expected no match for invalid JSON")
	}
}

func TestCompareFilter_TransformLine_Unchanged(t *testing.T) {
	f, _ := NewCompareFilter("x", OpNe, 0)
	line := `{"x":5}`
	if f.TransformLine(line) != line {
		t.Error("expected TransformLine to return line unchanged")
	}
}
