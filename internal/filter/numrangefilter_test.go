package filter

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func makeNumLine(field string, value interface{}) string {
	b, _ := json.Marshal(map[string]interface{}{field: value})
	return string(b)
}

func TestNewNumRangeFilter_Valid(t *testing.T) {
	f, err := NewNumRangeFilter("latency", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field != "latency" {
		t.Errorf("expected field latency, got %s", f.Field)
	}
}

func TestNewNumRangeFilter_EmptyField(t *testing.T) {
	_, err := NewNumRangeFilter("", 0, 100)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewNumRangeFilter_MinGtMax(t *testing.T) {
	_, err := NewNumRangeFilter("score", 50, 10)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNumRangeFilter_MatchesLine_InRange(t *testing.T) {
	f, _ := NewNumRangeFilter("ms", 10, 200)
	for _, v := range []float64{10, 100, 200} {
		line := makeNumLine("ms", v)
		if !f.MatchesLine(line) {
			t.Errorf("expected match for value %v", v)
		}
	}
}

func TestNumRangeFilter_MatchesLine_OutOfRange(t *testing.T) {
	f, _ := NewNumRangeFilter("ms", 10, 200)
	for _, v := range []float64{9.9, 200.1, -1} {
		line := makeNumLine("ms", v)
		if f.MatchesLine(line) {
			t.Errorf("expected no match for value %v", v)
		}
	}
}

func TestNumRangeFilter_MatchesLine_MissingField(t *testing.T) {
	f, _ := NewNumRangeFilter("ms", 0, 100)
	line := `{"other":42}`
	if f.MatchesLine(line) {
		t.Error("expected no match for missing field")
	}
}

func TestNumRangeFilter_MatchesLine_NonNumeric(t *testing.T) {
	f, _ := NewNumRangeFilter("ms", 0, 100)
	line := `{"ms":"fast"}`
	if f.MatchesLine(line) {
		t.Error("expected no match for non-numeric value")
	}
}

func TestNumRangeFilter_MatchesLine_InvalidJSON(t *testing.T) {
	f, _ := NewNumRangeFilter("ms", 0, 100)
	if f.MatchesLine("not-json") {
		t.Error("expected no match for invalid JSON")
	}
}

func TestNumRangeFilter_OpenBounds(t *testing.T) {
	f, _ := NewNumRangeFilter("score", math.Inf(-1), math.Inf(1))
	line := makeNumLine("score", 9999.9)
	if !f.MatchesLine(line) {
		t.Error("expected match with open bounds")
	}
}

func TestParseNumRangeFlag_Valid(t *testing.T) {
	f, err := ParseNumRangeFlag("latency:10:500")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Min != 10 || f.Max != 500 {
		t.Errorf("unexpected bounds: %v %v", f.Min, f.Max)
	}
}

func TestParseNumRangeFlag_OpenMin(t *testing.T) {
	f, err := ParseNumRangeFlag("score:*:100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !math.IsInf(f.Min, -1) {
		t.Errorf("expected -Inf min, got %v", f.Min)
	}
}

func TestParseNumRangeFlag_InvalidFormat(t *testing.T) {
	_, err := ParseNumRangeFlag("latency:10")
	if err == nil {
		t.Fatal("expected error for bad format")
	}
}

func TestParseNumRangeFlag_BadMin(t *testing.T) {
	_, err := ParseNumRangeFlag(fmt.Sprintf("ms:abc:100"))
	if err == nil {
		t.Fatal("expected error for non-numeric min")
	}
}
