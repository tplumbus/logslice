package filter

import (
	"testing"
)

func TestNewHeadTailFilter_Valid(t *testing.T) {
	for _, mode := range []string{"head", "tail"} {
		_, err := NewHeadTailFilter(mode, 5)
		if err != nil {
			t.Errorf("mode=%s: unexpected error: %v", mode, err)
		}
	}
}

func TestNewHeadTailFilter_InvalidMode(t *testing.T) {
	_, err := NewHeadTailFilter("middle", 5)
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestNewHeadTailFilter_ZeroN(t *testing.T) {
	_, err := NewHeadTailFilter("head", 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestNewHeadTailFilter_NegativeN(t *testing.T) {
	_, err := NewHeadTailFilter("tail", -3)
	if err == nil {
		t.Fatal("expected error for negative n")
	}
}

func TestHeadFilter_PassesFirstN(t *testing.T) {
	f, _ := NewHeadTailFilter("head", 3)
	lines := []string{"a", "b", "c", "d", "e"}
	var got []string
	for _, l := range lines {
		if f.MatchesLine(l) {
			got = append(got, l)
		}
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	for i, want := range []string{"a", "b", "c"} {
		if got[i] != want {
			t.Errorf("line %d: got %q, want %q", i, got[i], want)
		}
	}
}

func TestHeadFilter_FlushReturnsNil(t *testing.T) {
	f, _ := NewHeadTailFilter("head", 2)
	f.MatchesLine("x")
	if f.Flush() != nil {
		t.Fatal("expected nil flush for head mode")
	}
}

func TestTailFilter_MatchesLineAlwaysFalse(t *testing.T) {
	f, _ := NewHeadTailFilter("tail", 3)
	for _, l := range []string{"a", "b", "c", "d"} {
		if f.MatchesLine(l) {
			t.Errorf("tail MatchesLine should always return false, got true for %q", l)
		}
	}
}

func TestTailFilter_FlushLastN(t *testing.T) {
	f, _ := NewHeadTailFilter("tail", 3)
	for _, l := range []string{"a", "b", "c", "d", "e"} {
		f.MatchesLine(l)
	}
	out := f.Flush()
	want := []string{"c", "d", "e"}
	if len(out) != len(want) {
		t.Fatalf("expected %d lines, got %d", len(want), len(out))
	}
	for i := range want {
		if out[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, out[i], want[i])
		}
	}
}

func TestTailFilter_FewerThanN(t *testing.T) {
	f, _ := NewHeadTailFilter("tail", 10)
	for _, l := range []string{"x", "y"} {
		f.MatchesLine(l)
	}
	out := f.Flush()
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}
