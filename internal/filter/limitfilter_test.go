package filter

import (
	"testing"
)

func TestNewLimitFilter_Valid(t *testing.T) {
	f, err := NewLimitFilter(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewLimitFilter_ZeroInvalid(t *testing.T) {
	_, err := NewLimitFilter(0)
	if err == nil {
		t.Fatal("expected error for limit=0")
	}
}

func TestNewLimitFilter_NegativeInvalid(t *testing.T) {
	_, err := NewLimitFilter(-3)
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestLimitFilter_MatchesLine_UnderLimit(t *testing.T) {
	f, _ := NewLimitFilter(3)
	for i := 0; i < 3; i++ {
		if !f.MatchesLine(`{"msg":"hello"}`) {
			t.Fatalf("expected true on iteration %d", i)
		}
	}
}

func TestLimitFilter_MatchesLine_AtLimit(t *testing.T) {
	f, _ := NewLimitFilter(2)
	f.MatchesLine(`{"a":"1"}`)
	f.MatchesLine(`{"a":"2"}`)
	if f.MatchesLine(`{"a":"3"}`) {
		t.Fatal("expected false after limit reached")
	}
}

func TestLimitFilter_Reset(t *testing.T) {
	f, _ := NewLimitFilter(1)
	f.MatchesLine(`{"a":"1"}`)
	if f.MatchesLine(`{"a":"2"}`) {
		t.Fatal("expected false after limit")
	}
	f.Reset()
	if !f.MatchesLine(`{"a":"3"}`) {
		t.Fatal("expected true after reset")
	}
}

func TestLimitFilter_MatchesLine_LimitOne(t *testing.T) {
	f, _ := NewLimitFilter(1)
	if !f.MatchesLine(`{"x":"y"}`) {
		t.Fatal("first line should match")
	}
	for i := 0; i < 5; i++ {
		if f.MatchesLine(`{"x":"y"}`) {
			t.Fatalf("line %d should not match after limit=1", i+2)
		}
	}
}
