package filter

import (
	"testing"
)

func TestNewSliceFilter_Valid(t *testing.T) {
	f, err := NewSliceFilter(2, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewSliceFilter_NegativeStart(t *testing.T) {
	_, err := NewSliceFilter(-1, 5)
	if err == nil {
		t.Fatal("expected error for negative start")
	}
}

func TestNewSliceFilter_EndBeforeStart(t *testing.T) {
	_, err := NewSliceFilter(5, 3)
	if err == nil {
		t.Fatal("expected error when end <= start")
	}
}

func TestNewSliceFilter_OpenEnd(t *testing.T) {
	f, err := NewSliceFilter(3, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestSliceFilter_MatchesLine_Window(t *testing.T) {
	f, _ := NewSliceFilter(1, 3)
	lines := []string{`{"i":0}`, `{"i":1}`, `{"i":2}`, `{"i":3}`}
	expected := []bool{false, true, true, false}
	for idx, line := range lines {
		got := f.MatchesLine(line)
		if got != expected[idx] {
			t.Errorf("line %d: expected %v, got %v", idx, expected[idx], got)
		}
	}
}

func TestSliceFilter_MatchesLine_OpenEnd(t *testing.T) {
	f, _ := NewSliceFilter(2, -1)
	results := make([]bool, 5)
	for i := range results {
		results[i] = f.MatchesLine(`{"x":1}`)
	}
	for i, v := range results {
		want := i >= 2
		if v != want {
			t.Errorf("index %d: want %v got %v", i, want, v)
		}
	}
}

func TestParseSliceFlag_Valid(t *testing.T) {
	start, end, err := ParseSliceFlag("2:10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if start != 2 || end != 10 {
		t.Errorf("expected 2:10, got %d:%d", start, end)
	}
}

func TestParseSliceFlag_OpenEnd(t *testing.T) {
	start, end, err := ParseSliceFlag("5:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if start != 5 || end != -1 {
		t.Errorf("expected 5:-1, got %d:%d", start, end)
	}
}

func TestParseSliceFlag_MissingColon(t *testing.T) {
	_, _, err := ParseSliceFlag("5")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseSliceFlag_BadStart(t *testing.T) {
	_, _, err := ParseSliceFlag("abc:10")
	if err == nil {
		t.Fatal("expected error for non-numeric start")
	}
}
